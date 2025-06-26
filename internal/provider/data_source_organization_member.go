package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mzglinski/go-sentry/v2/sentry"
	"github.com/mzglinski/terraform-provider-sentry/internal/diagutils"
)

type OrganizationMemberDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	Organization types.String `tfsdk:"organization"`
	UserId       types.String `tfsdk:"user_id"`
	Email        types.String `tfsdk:"email"`
	Role         types.String `tfsdk:"role"`
}

func (m *OrganizationMemberDataSourceModel) Fill(ctx context.Context, member *sentry.OrganizationMember) (diags diag.Diagnostics) {
	m.Id = types.StringValue(member.ID)
	if member.User.ID != "" {
		m.UserId = types.StringValue(member.User.ID)
	} else {
		m.UserId = types.StringNull()
	}
	m.Email = types.StringValue(member.Email)
	m.Role = types.StringValue(member.OrgRole)
	return
}

var _ datasource.DataSource = &OrganizationMemberDataSource{}
var _ datasource.DataSourceWithConfigure = &OrganizationMemberDataSource{}

func NewOrganizationMemberDataSource() datasource.DataSource {
	return &OrganizationMemberDataSource{}
}

type OrganizationMemberDataSource struct {
	baseDataSource
}

func (d *OrganizationMemberDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_member"
}

func (d *OrganizationMemberDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Retrieve an organization member by email.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of this resource.",
				Computed:            true,
			},
			"organization": DataSourceOrganizationAttribute(),
			"user_id": schema.StringAttribute{
				MarkdownDescription: "The user ID of the organization member.",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email of the organization member.",
				Required:            true,
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "This is the role of the organization member.",
				Computed:            true,
			},
		},
	}
}

func (d *OrganizationMemberDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationMemberDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var foundMember *sentry.OrganizationMember
	params := &sentry.ListCursorParams{}

out:
	for {
		members, sentryResp, err := d.client.OrganizationMembers.List(ctx, data.Organization.ValueString(), params)
		if err != nil {
			resp.Diagnostics.Append(diagutils.NewClientError("read", err))
			return
		}

		for _, member := range members {
			if member.Email == data.Email.ValueString() {
				foundMember = member
				break out
			}
		}

		if sentryResp.Cursor == "" {
			break
		}
		params.Cursor = sentryResp.Cursor
	}

	if foundMember == nil {
		resp.Diagnostics.AddError("Not found", "No matching organization member found")
		return
	}

	resp.Diagnostics.Append(data.Fill(ctx, foundMember)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
