package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/twmb/franz-go/pkg/kadm"
	"strconv"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &BrokerDataSource{}

func NewBrokerDataSource() datasource.DataSource {
	return &BrokerDataSource{}
}

// BrokerDataSource defines the data source implementation.
type BrokerDataSource struct {
	client *kadm.Client
}

// BrokerDataSourceModel describes the data source data model.
type BrokerDataSourceModel struct {
	Brokers map[string]BrokerDataModel `tfsdk:"brokers"`
}

// BrokerDataModel describes the data source data model.
type BrokerDataModel struct {
	NodeID types.Int64  `tfsdk:"nodeid"`
	Host   types.String `tfsdk:"host"`
	Port   types.Int64  `tfsdk:"port"`
	RackID types.String `tfsdk:"rackid"`
}

func (d *BrokerDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_brokers"
}

func (d *BrokerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Brokers data source",
		Attributes: map[string]schema.Attribute{
			"brokers": schema.MapNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"nodeid": schema.Int64Attribute{
							MarkdownDescription: "nodeid",
							Computed:            true,
						},
						"host": schema.StringAttribute{
							MarkdownDescription: "host",
							Computed:            true,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "port",
							Computed:            true,
						},
						"rackid": schema.StringAttribute{
							MarkdownDescription: "host",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *BrokerDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kadm.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *kadm.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *BrokerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BrokerDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	metadata, err := d.client.BrokerMetadata(ctx)
	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("unable to retrieve broker metadata: %v", err))
	} else {
		brokers := make(map[string]BrokerDataModel, len(metadata.Brokers.NodeIDs()))
		for _, broker := range metadata.Brokers {
			brokers[strconv.Itoa(int(broker.NodeID))] = BrokerDataModel{
				NodeID: types.Int64Value(int64(broker.NodeID)),
				Host:   types.StringValue(broker.Host),
				Port:   types.Int64Value(int64(broker.Port)),
				RackID: types.StringValue(*broker.Rack),
			}
		}
		data.Brokers = brokers
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
