package devhelpermcp

import (
	"errors"
	"fmt"
	"time"

	"google.golang.org/adk/tool"
	"google.golang.org/protobuf/types/known/durationpb"

	artifactregistry "cloud.google.com/go/artifactregistry/apiv1"
	artifactregistrypb "cloud.google.com/go/artifactregistry/apiv1/artifactregistrypb"
)

var (
	moduleCleanup = 7 * 24 * time.Hour
)

type createGoRepositoryInput struct {
}

func (i createGoRepositoryInput) validate() error {
	var errs []error
	return errors.Join(errs...)
}

type createGoRepositoryOutput struct {
	ProjectID      string `json:"project_id" jsonschema:"The project ID of the created Go module repository.`
	Region         string `json:"region" jsonschema:"The region of the created Go module repository.`
	RepositoryID   string `json:"repository_id" jsonschema:"The repository ID of the created Go module repository.`
	RepositoryName string `json:"repository_name" jsonschema:"The full resource name of the created Go module repository.`
	RegistryUri    string `json:"registry_uri" jsonschema:"The URI of the created Go module repository.`
}

func (d *dh) createGoRepositoryTool(ctx tool.Context, input createGoRepositoryInput) (createGoRepositoryOutput, error) {
	var zeroOutput createGoRepositoryOutput
	if err := input.validate(); err != nil {
		return zeroOutput, err
	}

	now := time.Now()
	repositoryID := fmt.Sprintf("%s-go-%s-%d", d.ServiceName, now.Format("20060102"), now.Unix())

	// TODO: possibly scrub old repositories to preserve quota?

	client, err := artifactregistry.NewClient(ctx)
	if err != nil {
		return zeroOutput, fmt.Errorf("artifactregistry.NewClient: %w", err)
	}
	defer client.Close()

	req := &artifactregistrypb.CreateRepositoryRequest{
		Parent:       fmt.Sprintf("projects/%s/locations/%s", d.ProjectID, d.Region),
		RepositoryId: repositoryID,
		Repository: &artifactregistrypb.Repository{
			Format:      artifactregistrypb.Repository_GO,
			Description: "Ephemeral Go module repository",
			VulnerabilityScanningConfig: &artifactregistrypb.Repository_VulnerabilityScanningConfig{
				EnablementConfig: artifactregistrypb.Repository_VulnerabilityScanningConfig_INHERITED,
			},
			CleanupPolicies: map[string]*artifactregistrypb.CleanupPolicy{
				"delete-after-one-day": &artifactregistrypb.CleanupPolicy{
					Action: artifactregistrypb.CleanupPolicy_DELETE,
					ConditionType: &artifactregistrypb.CleanupPolicy_Condition{
						Condition: &artifactregistrypb.CleanupPolicyCondition{
							OlderThan: durationpb.New(moduleCleanup),
						},
					},
				},
			},
		},
	}

	op, err := client.CreateRepository(ctx, req)
	if err != nil {
		return zeroOutput, fmt.Errorf("client.CreateRepository failed: %w", err)
	}

	// Wait for the LRO to complete and retrieve the result.
	repo, err := op.Wait(ctx)
	if err != nil {
		return zeroOutput, fmt.Errorf("op.Wait failed: %w", err)
	}

	return createGoRepositoryOutput{
		ProjectID:      d.ProjectID,
		Region:         d.Region,
		RepositoryID:   repositoryID,
		RepositoryName: repo.Name,
		RegistryUri:    repo.RegistryUri,
	}, nil
}
