package db

import (
	"fmt"

	"github.com/gbolo/vsummary/common"
)

// InsertVMs inserts a vm into database
func (b *Backend) InsertVMs(vms []common.Vm) (err error) {

	// exit if there is no database connection
	err = b.checkDB()
	if err != nil {
		return
	}

	// begin a transaction, set all related objects to absent
	tx := b.db.MustBegin()
	// TODO: improve this
	tx.MustExec("UPDATE vm SET present = 0 WHERE present = 1 AND vcenter_id=?", vms[0].VcenterId)
	var rowsAffected int64 = 0

	for _, vm := range vms {
		// fill in missing data
		// folder may not exist
		if vm.FolderMoref != "" && vm.FolderMoref != "vapp" {
			vm.FolderId = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.FolderMoref))
		}

		// vapps may not exist
		if vm.VappMoref != "none" {
			vm.VappId = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.VappMoref))
		}
		// resourcepool may not exist
		if vm.ResourcePoolMoref != "" {
			vm.ResourcePoolId = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.ResourcePoolId))
		}

		// Fill in some required Ids
		vm.Id = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.Moref))
		vm.EsxiId = common.GetMD5Hash(fmt.Sprintf("%s%s", vm.VcenterId, vm.EsxiMoref))

		// Store the user record in the DB
		res, err := tx.NamedExec(insertVm, &vm)

		if err != nil {
			break
		}

		// tally up rows affected for logging
		numRowsAffected, err := res.RowsAffected()
		if err != nil {
			break
		}
		rowsAffected = rowsAffected + numRowsAffected

	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		log.Errorf("failed to commit transaction to database: %s", err)
	}

	log.Debugf("total combined affected rows: %d", rowsAffected)

	return

}
