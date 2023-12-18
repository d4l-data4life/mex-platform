package index

import (
	"context"
	"fmt"

	"github.com/d4l-data4life/mex/mex/shared/constants"
	"github.com/d4l-data4life/mex/mex/shared/known/jobspb"
	"github.com/d4l-data4life/mex/mex/shared/known/statuspb"
	L "github.com/d4l-data4life/mex/mex/shared/log"
	"github.com/d4l-data4life/mex/mex/shared/solr"
	"github.com/d4l-data4life/mex/mex/shared/utils"
)

// runRecreateCollectionJob starts a job that completely drops and rebuilds the Solr collection used by MEx,
// using the newest configs. This ensures a complete removal of the old index.
func (svc *Service) runRecreateCollectionJob(configHash string) error {
	ctx := context.Background()
	lock, err := svc.JobService.AcquireLock(ctx, SvcResourceName)
	if err != nil {
		return fmt.Errorf("recreate: failed to acquire index lock; other job might be running: %s", err.Error())
	}

	job, err := svc.JobService.CreateJob(ctx, &jobspb.CreateJobRequest{
		Title: fmt.Sprintf("Recreate Solr collection due to config change (hash %q)", configHash),
	})
	if err != nil {
		return fmt.Errorf("failure creating job:  %s", err.Error())
	}

	// Create new independent context for the job copying the relevant values.
	// (The request's ctx will go out of scope before the job is done, so we cannot use it directly or as a parent context.)
	ctxJob := constants.NewContextWithValues(ctx, job.JobId)

	go func(ctx context.Context) {
		svc.Log.Info(ctx, L.Messagef("collection recreation: job started (%s)", job.JobId), L.Phase("job"))
		svc.TelemetryService.SetStatus(statuspb.Color_AMBER, configHash)

		svc.JobService.SetStatusRunning(ctx, job.JobId) //nolint:errcheck

		var statusColor statuspb.Color
		err = svc.doCollectionRecreate(ctx, svc.TelemetryService)
		if err != nil {
			svc.Log.Error(ctx, L.Message(err.Error()))
			statusColor = statuspb.Color_RED

			_, err := svc.JobService.SetError(ctx, &jobspb.SetJobErrorRequest{Error: err.Error(), JobId: job.JobId})
			if err != nil {
				svc.Log.Warn(ctx, L.Messagef("could not set job error: %s", err.Error()))
			}
		} else {
			statusColor = statuspb.Color_GREEN
		}
		svc.Log.Info(ctx, L.Messagef("collection recreation: job done (%s)", job.JobId), L.Phase("job"))
		svc.TelemetryService.Done()

		svc.JobService.SetStatusDone(ctx, job.JobId)           //nolint:errcheck
		svc.JobService.ReleaseLock(ctx, SvcResourceName, lock) //nolint:errcheck

		// Make sure we set the status after releasing the index resource.
		svc.TelemetryService.SetStatus(statusColor, configHash)
	}(ctxJob)

	return nil
}

// doCollectionRecreate carries out the individual steps needed for a complete collection rebuild
func (svc *Service) doCollectionRecreate(ctx context.Context, progressor utils.Progressor) error {
	if progressor == nil {
		progressor = &utils.NopProgressor{}
	}

	progressor.Progress("re-create collection", "started")

	svc.Log.Info(ctx, L.Messagef("recreating and reindexing Solr collection '%s'", svc.CollectionName))
	collections, err := svc.Solr.GetCollections(ctx)
	if err != nil {
		svc.Log.Error(ctx, L.Message("Failed to get Solr collections"))
		return err
	}
	if utils.Contains(collections, svc.CollectionName) {
		err := svc.Solr.DeleteCollection(ctx, svc.CollectionName)
		if err != nil {
			svc.Log.Error(ctx, L.Messagef("Failed to delete Solr collection '%s'", svc.CollectionName))
			return err
		}
	}

	err = svc.Solr.CreateCollection(ctx, svc.CollectionName, solr.DefaultSolrConfigSet, svc.ReplicationFactor)
	if err != nil {
		svc.Log.Error(ctx, L.Messagef("Failed to create Solr collection '%s'", svc.CollectionName))
		return err
	}
	progressor.Progress("creating collection", "done")

	err = svc.doSchemaRebuild(ctx, progressor)
	if err != nil {
		svc.Log.Error(ctx, L.Message("Failed to rebuild Solr schema"))
		return err
	}
	progressor.Progress("rebuilding schema", "done")

	err = svc.DoIndexUpdate(ctx, progressor)
	if err != nil {
		svc.Log.Error(ctx, L.Message("Failed to update Solr index"))
		return err
	}
	progressor.Progress("reindexing", "done")

	svc.Log.Info(ctx, L.Messagef("Solr collection '%s' was recreated and reindexed", svc.CollectionName))

	return nil
}
