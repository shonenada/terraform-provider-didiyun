package didiyun

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	ddy "github.com/shonenada/didiyun-go"
	job "github.com/shonenada/didiyun-go/job"
)

const TimeOut = 2 * time.Minute

func WaitForJob(client *ddy.Client, regionId string, jobUuid string) error {
	return resource.Retry(TimeOut, func() *resource.RetryError {
		jobs, err := client.Job().GetResult(&job.ResultRequest{
			RegionId: regionId,
			Uuids:    jobUuid,
		})
		if err != nil {
			return resource.RetryableError(fmt.Errorf("Failed to get job: %v", err))
		}

		job := (*jobs)[0]

		if job.Progress < 100 {
			return resource.RetryableError(fmt.Errorf("Wait for job"))
		}

		if !job.IsDone {
			return resource.RetryableError(fmt.Errorf("Wait for job"))
		}

		if !job.IsSuccess {
			return resource.NonRetryableError(fmt.Errorf("Failed to execute job: %v", job.Result))
		}

		return nil
	})
}
