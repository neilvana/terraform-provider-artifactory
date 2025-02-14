package local

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jfrog/terraform-provider-artifactory/v6/pkg/artifactory/resource/repository"
	"github.com/jfrog/terraform-provider-shared/util"
	"github.com/jfrog/terraform-provider-shared/validator"
)

func ResourceArtifactoryLocalRpmRepository() *schema.Resource {
	const packageType = "rpm"

	var rpmLocalSchema = util.MergeSchema(BaseLocalRepoSchema, map[string]*schema.Schema{
		"yum_root_depth": {
			Type:             schema.TypeInt,
			Optional:         true,
			Default:          0,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
			Description: "The depth, relative to the repository's root folder, where RPM metadata is created. " +
				"This is useful when your repository contains multiple RPM repositories under parallel hierarchies. " +
				"For example, if your RPMs are stored under 'fedora/linux/$releasever/$basearch', specify a depth of 4.",
		},
		"calculate_yum_metadata": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"enable_file_lists_indexing": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"yum_group_file_names": {
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "",
			ValidateDiagFunc: validator.CommaSeperatedList,
			Description: "A comma separated list of XML file names containing RPM group component definitions. Artifactory includes " +
				"the group definitions as part of the calculated RPM metadata, as well as automatically generating a " +
				"gzipped version of the group files, if required.",
		},
		"primary_keypair_ref": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			Description:      "Primary keypair used to sign artifacts.",
		},
		"secondary_keypair_ref": {
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
			Description:      "Secondary keypair used to sign artifacts.",
		},
	}, repository.RepoLayoutRefSchema("local", packageType))

	type RpmLocalRepositoryParams struct {
		LocalRepositoryBaseParams
		RootDepth               int    `hcl:"yum_root_depth" json:"yumRootDepth"`
		CalculateYumMetadata    bool   `hcl:"calculate_yum_metadata" json:"calculateYumMetadata"`
		EnableFileListsIndexing bool   `hcl:"enable_file_lists_indexing" json:"enableFileListsIndexing"`
		GroupFileNames          string `hcl:"yum_group_file_names" json:"yumGroupFileNames"`
		PrimaryKeyPairRef       string `json:"primaryKeyPairRef"`
		SecondaryKeyPairRef     string `json:"secondaryKeyPairRef"`
	}

	unPackLocalRpmRepository := func(data *schema.ResourceData) (interface{}, string, error) {
		d := &util.ResourceData{ResourceData: data}
		repo := RpmLocalRepositoryParams{
			LocalRepositoryBaseParams: UnpackBaseRepo("local", data, "rpm"),
			RootDepth:                 d.GetInt("yum_root_depth", false),
			CalculateYumMetadata:      d.GetBool("calculate_yum_metadata", false),
			EnableFileListsIndexing:   d.GetBool("enable_file_lists_indexing", false),
			GroupFileNames:            d.GetString("yum_group_file_names", false),
			PrimaryKeyPairRef:         d.GetString("primary_keypair_ref", false),
			SecondaryKeyPairRef:       d.GetString("secondary_keypair_ref", false),
		}

		return repo, repo.Id(), nil
	}

	return repository.MkResourceSchema(rpmLocalSchema, repository.DefaultPacker(rpmLocalSchema), unPackLocalRpmRepository, func() interface{} {
		return &RpmLocalRepositoryParams{
			LocalRepositoryBaseParams: LocalRepositoryBaseParams{
				PackageType: packageType,
				Rclass:      "local",
			},
			RootDepth:               0,
			CalculateYumMetadata:    false,
			EnableFileListsIndexing: false,
			GroupFileNames:          "",
		}
	})
}
