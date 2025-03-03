package cli

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/andydunstall/piko/cli/agent"
	"github.com/andydunstall/piko/cli/bench"
	"github.com/andydunstall/piko/cli/forward"
	"github.com/andydunstall/piko/cli/server"
	"github.com/andydunstall/piko/cli/test"
	"github.com/andydunstall/piko/pkg/build"
)

func NewCommand() *cobra.Command {
	buildModules := build.Module
	shouldBuildServer := buildModules == "server" || buildModules == "all"
	shouldBuildClient := buildModules == "client" || buildModules == "all"

	var commandDescriptionBuilder strings.Builder

	commandDescriptionBuilder.WriteString(`
	AMP Ingress is a reverse proxy that allows you to 
	expose an endpoint that isn't publicly routable to AmpID service.`)

	if shouldBuildServer {
		commandDescriptionBuilder.WriteString(`
The Piko server is responsible for routing incoming proxy requests to upstream
services. Upstream services open outbound-connections to the server and
register endpoints. Piko will then route incoming requests to the appropriate
upstream service via the upstreams outbound-only connection.

The server may be hosted as a cluster of nodes.

Start a server node with:

  $ piko server

You can also inspect the status of the server using:

  $ piko server status
		`)
	}

	if shouldBuildClient {
		commandDescriptionBuilder.WriteString(`
To register an upstream service, use the Ingress agent. The agent is a lightweight
proxy that runs alongside your services. It connects to the AMP Ingress server,
registers the configured endpoints, then forwards incoming requests to your
services.

Such as to register endpoint 'my-endpoint' to forward HTTP requests to your
service at 'localhost:3000':

  $ amp_ingress agent http my-endpoint 3000

You can also forward raw TCP using:

  $ amp_ingress agent tcp my-endpoint 3000

To forward a local TCP port to an upstream endpoint, use 'amp_ingress forward'.
This listens for TCP connections on the configured local port and forwards them
to an upstream listener via AMP. Such as to forward port 3000 to endpoint
'my-endpoint':

  $ amp_ingress forward tcp 3000 my-endpoint`)
	}

	cmd := &cobra.Command{
		Use:          "amp_ingress [command] (flags)",
		SilenceUsage: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		Long:    commandDescriptionBuilder.String(),
		Version: build.Version,
	}

	if shouldBuildServer {
		cmd.AddCommand(server.NewCommand())
		cmd.AddCommand(bench.NewCommand())
		cmd.AddCommand(test.NewCommand())
	}
	if shouldBuildClient {
		cmd.AddCommand(agent.NewCommand())
		cmd.AddCommand(forward.NewCommand())
	}

	return cmd
}

func init() {
	cobra.EnableCommandSorting = false
}
