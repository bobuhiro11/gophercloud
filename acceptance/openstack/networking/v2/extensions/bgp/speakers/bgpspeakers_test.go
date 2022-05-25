package speakers

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	ap "github.com/gophercloud/gophercloud/acceptance/openstack/networking/v2/extensions/bgp/peers"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgp/peers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgp/speakers"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestBGPSpeakerCRUD(t *testing.T) {
	clients.RequireAdmin(t)
	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	// Create a BGP Speaker
	bgpSpeaker, err := CreateBGPSpeaker(t, client)
	th.AssertNoErr(t, err)

	// Create a BGP Peer
	bgpPeer, err := ap.CreateBGPPeer(t, client)
	th.AssertNoErr(t, err)

	// List BGP Speakers
	allPages, err := speakers.List(client).AllPages()
	th.AssertNoErr(t, err)
	allSpeakers, err := speakers.ExtractBGPSpeakers(allPages)
	th.AssertNoErr(t, err)

	t.Logf("Retrieved BGP Speakers")
	tools.PrintResource(t, allSpeakers)
	th.AssertIntGreaterOrEqual(t, len(allSpeakers), 1)

	// Update BGP Speaker
	opts := speakers.UpdateOpts{
		Name:                          tools.RandomString("TESTACC-BGPSPEAKER-", 10),
		AdvertiseTenantNetworks:       false,
		AdvertiseFloatingIPHostRoutes: true,
	}
	speakerUpdated, err := speakers.Update(client, bgpSpeaker.ID, opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, speakerUpdated.Name, opts.Name)
	t.Logf("Updated the BGP Speaker, name set from %s to %s", bgpSpeaker.Name, speakerUpdated.Name)

	// Get a BGP Speaker
	bgpSpeakerGot, err := speakers.Get(client, bgpSpeaker.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, bgpSpeaker.ID, bgpSpeakerGot.ID)
	th.AssertEquals(t, opts.Name, bgpSpeakerGot.Name)

	// AddBGPPeer
	addBGPPeerOpts := speakers.AddBGPPeerOpts{BGPPeerID: bgpPeer.ID}
	_, err = speakers.AddBGPPeer(client, bgpSpeaker.ID, addBGPPeerOpts).Extract()
	th.AssertNoErr(t, err)
	speakerGot, err := speakers.Get(client, bgpSpeaker.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, bgpPeer.ID, speakerGot.Peers[0])
	t.Logf("Successfully added BGP Peer %s to BGP Speaker %s", bgpPeer.Name, speakerUpdated.Name)

	// RemoveBGPPeer
	removeBGPPeerOpts := speakers.RemoveBGPPeerOpts{BGPPeerID: bgpPeer.ID}
	err = speakers.RemoveBGPPeer(client, bgpSpeaker.ID, removeBGPPeerOpts).ExtractErr()
	th.AssertNoErr(t, err)
	speakerGot, err = speakers.Get(client, bgpSpeaker.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, len(speakerGot.Networks), 0)
	t.Logf("Successfully removed BGP Peer %s to BGP Speaker %s", bgpPeer.Name, speakerUpdated.Name)

	// Delete a BGP Peer
	t.Logf("Delete the BGP Peer %s", bgpPeer.Name)
	err = peers.Delete(client, bgpPeer.ID).ExtractErr()
	th.AssertNoErr(t, err)

	// Delete a BGP Speaker
	t.Logf("Delete the BGP Speaker %s", speakerUpdated.Name)
	err = speakers.Delete(client, bgpSpeaker.ID).ExtractErr()
	th.AssertNoErr(t, err)
}