package utils

import (
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
)

// Copied from mgm hooks
func CallToBeforeCreateHooks(model mgm.Model) error {
	if hook, ok := model.(mgm.CreatingHook); ok {
		if err := hook.Creating(); err != nil {
			return err
		}
	}

	if hook, ok := model.(mgm.SavingHook); ok {
		if err := hook.Saving(); err != nil {
			return err
		}
	}

	return nil
}

func CallToBeforeUpdateHooks(model mgm.Model) error {
	if hook, ok := model.(mgm.UpdatingHook); ok {
		if err := hook.Updating(); err != nil {
			return err
		}
	}

	if hook, ok := model.(mgm.SavingHook); ok {
		if err := hook.Saving(); err != nil {
			return err
		}
	}

	return nil
}

func CallToAfterCreateHooks(model mgm.Model) error {
	if hook, ok := model.(mgm.CreatedHook); ok {
		if err := hook.Created(); err != nil {
			return err
		}
	}

	if hook, ok := model.(mgm.SavedHook); ok {
		if err := hook.Saved(); err != nil {
			return err
		}
	}

	return nil
}

func CallToAfterUpdateHooks(updateResult *mongo.UpdateResult, model mgm.Model) error {
	if hook, ok := model.(mgm.UpdatedHook); ok {
		if err := hook.Updated(updateResult); err != nil {
			return err
		}
	}

	if hook, ok := model.(mgm.SavedHook); ok {
		if err := hook.Saved(); err != nil {
			return err
		}
	}

	return nil
}

func CallToBeforeDeleteHooks(model mgm.Model) error {
	if hook, ok := model.(mgm.DeletingHook); ok {
		if err := hook.Deleting(); err != nil {
			return err
		}
	}

	return nil
}

func CallToAfterDeleteHooks(deleteResult *mongo.DeleteResult, model mgm.Model) error {
	if hook, ok := model.(mgm.DeletedHook); ok {
		if err := hook.Deleted(deleteResult); err != nil {
			return err
		}
	}

	return nil
}