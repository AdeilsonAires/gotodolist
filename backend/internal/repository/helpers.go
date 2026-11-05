package repository

import (
	"errors"

	"github.com/AdeilsonAires/gotodolist/backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// checkAffected traduz o resultado de um Update/Delete em erro de domínio.
func checkAffected(result *gorm.DB) error {
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// scoped devolve uma query de itens já restrita à lista E ao dono dela.
// Recebe o *gorm.DB (normal ou de transação) para funcionar nos dois casos.
func scoped(db *gorm.DB, listID, userID uint) *gorm.DB {
	ownedList := db.Model(&domain.List{}).
		Select("id").
		Where("id = ? AND user_id = ?", listID, userID)

	return db.Model(&domain.Item{}).
		Where("list_id IN (?)", ownedList)
}

// lockList confirma que a lista é do usuário e a trava (SELECT ... FOR UPDATE)
// até o fim da transação, evitando que operações simultâneas calculem o
// status da lista com dados desatualizados.
func lockList(tx *gorm.DB, listID, userID uint) error {
	var list domain.List
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Select("id").
		Where("id = ? AND user_id = ?", listID, userID).
		First(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrNotFound
	}
	return err
}

// syncListCompletion recalcula o status da lista a partir dos itens.
// Regra: concluída = tem pelo menos 1 item E todos estão concluídos.
// Deve ser chamada dentro da transação que alterou os itens.
func syncListCompletion(tx *gorm.DB, listID uint) error {
	var counts struct {
		Total int64
		Done  int64
	}
	err := tx.Model(&domain.Item{}).
		Select("COUNT(*) AS total, COUNT(*) FILTER (WHERE completed) AS done").
		Where("list_id = ?", listID).
		Scan(&counts).Error
	if err != nil {
		return err
	}

	completed := counts.Total > 0 && counts.Total == counts.Done

	return tx.Model(&domain.List{}).
		Where("id = ?", listID).
		Updates(map[string]any{
			"completed": completed,
			// mantém a data original se já estava concluída; zera se voltou a ativa
			"completed_at": gorm.Expr("CASE WHEN ? THEN COALESCE(completed_at, NOW()) ELSE NULL END", completed),
		}).Error
}