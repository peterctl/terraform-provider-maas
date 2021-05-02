package maas

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ionutbalutoiu/gomaasclient/client"
)

func dataSourceMaasSubnet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSubnetRead,

		Schema: map[string]*schema.Schema{
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
			},
			"vlan_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fabric": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"space": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_servers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rdns_mode": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceSubnetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	subnets, err := client.Subnets.Get()
	if err != nil {
		return diag.FromErr(err)
	}

	cidr := d.Get("cidr").(string)
	vlanId, vlanIdOk := d.GetOk("vlan_id")

	for _, subnet := range subnets {
		if cidr != subnet.CIDR {
			continue
		}
		if vlanIdOk {
			if vlanId.(int) != subnet.VLAN.ID {
				continue
			}
		}
		if err := d.Set("vid", subnet.VLAN.VID); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("fabric", subnet.VLAN.Fabric); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("name", subnet.Name); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("space", subnet.Space); err != nil {
			return diag.FromErr(err)
		}
		gatewayIp := subnet.GatewayIP.String()
		if gatewayIp == "<nil>" {
			gatewayIp = ""
		}
		if err := d.Set("gateway_ip", gatewayIp); err != nil {
			return diag.FromErr(err)
		}
		dnsServers := make([]string, len(subnet.DNSServers))
		for i, ip := range subnet.DNSServers {
			dnsServers[i] = ip.String()
		}
		if err := d.Set("dns_servers", dnsServers); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("rdns_mode", subnet.RDNSMode); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(fmt.Sprintf("%v", subnet.ID))
		return nil
	}

	return diag.FromErr(fmt.Errorf("could not find matching subnet"))
}
