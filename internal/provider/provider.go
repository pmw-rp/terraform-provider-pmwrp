// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl"
	"github.com/twmb/franz-go/pkg/sasl/scram"
	"github.com/twmb/tlscfg"
	"net"
	"time"
)

// Ensure PMWRPProvider satisfies various provider interfaces.
var _ provider.Provider = &PMWRPProvider{}

// PMWRPProvider defines the provider implementation.
type PMWRPProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// PMWRPProviderModel describes the provider data model.
type PMWRPProviderModel struct {
	Seed     types.String `tfsdk:"seed"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *PMWRPProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pmwrp"
	resp.Version = p.version
}

func (p *PMWRPProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"seed": schema.StringAttribute{
				MarkdownDescription: "seed broker address",
				Required:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "username",
				Required:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "password",
				Required:            true,
			},
		},
	}
}

func (p *PMWRPProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PMWRPProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Seed.IsNull() { /* ... */ }

	tflog.Info(ctx, "connecting to "+data.Seed.String())

	var mechanism sasl.Mechanism
	scramAuth := scram.Auth{
		User: data.Username.ValueString(),
		Pass: data.Password.ValueString(),
	}
	mechanism = scramAuth.AsSha256Mechanism()

	tlsCfg, err := tlscfg.New()

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("failed to create tls config: %v", err))
		//return nil, err
	}

	tlsDialer := &tls.Dialer{
		NetDialer: &net.Dialer{Timeout: 10 * time.Second},
		Config:    tlsCfg,
	}

	var client *kadm.Client
	{
		cl, err := kgo.NewClient(
			kgo.SeedBrokers(data.Seed.ValueString()),
			kgo.SASL(sasl.Mechanism(mechanism)),
			kgo.Dialer(tlsDialer.DialContext),
		)

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("unable to create admin client: %v", err))
			//return nil, err
		}
		client = kadm.NewClient(cl)
	}

	// Example client configuration for data sources and resources
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PMWRPProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *PMWRPProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewBrokerDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &PMWRPProvider{
			version: version,
		}
	}
}
