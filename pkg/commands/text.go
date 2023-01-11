package commands

import (
	"fmt"
	"io"
	"log"
	"os"

	"chainguard.dev/apko/pkg/build/types"
	"github.com/dominikbraun/graph"
	"github.com/spf13/cobra"
	"github.com/wolfi-dev/dag/pkg"
)

func cmdText() *cobra.Command {
	var dir, arch, t string
	var showDependents bool
	text := &cobra.Command{
		Use:   "text",
		Short: "Print a sorted list of downstream dependent packages",
		RunE: func(cmd *cobra.Command, args []string) error {
			arch := types.ParseArchitecture(arch).ToAPK()

			g, err := pkg.NewGraph(os.DirFS(dir), dir)
			if err != nil {
				return err
			}

			if len(args) == 0 {
				if showDependents {
					log.Print("warning: the 'show dependents' option has no effect without specifying one or more package names")
				}
			} else {
				// ensure all packages exist in the graph
				for _, arg := range args {
					if _, err := g.Graph.Vertex(arg); err == graph.ErrVertexNotFound {
						return fmt.Errorf("package %q not found in graph", arg)
					}
				}

				// determine if we're examining dependencies or dependents
				var subgraph *pkg.Graph
				if showDependents {
					leaves := args
					subgraph, err = g.SubgraphWithLeaves(leaves)
					if err != nil {
						return err
					}
				} else {
					roots := args
					subgraph, err = g.SubgraphWithRoots(roots)
					if err != nil {
						return err
					}
				}

				g = subgraph
			}

			return text(*g, arch, textType(t), os.Stdout)
		},
	}
	text.Flags().StringVarP(&dir, "dir", "d", ".", "directory to search for melange configs")
	text.Flags().StringVarP(&arch, "arch", "a", "x86_64", "architecture to build for")
	text.Flags().BoolVarP(&showDependents, "show-dependents", "D", false, "show packages that depend on these packages, instead of these packages' dependencies")
	text.Flags().StringVarP(&t, "type", "t", string(typeTarget), "What type of text to emit")
	return text
}

type textType string

const (
	typeTarget       = "target"
	typeMakefileLine = "makefile"
)

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}

func text(g pkg.Graph, arch string, t textType, w io.Writer) error {
	all, err := g.Sorted()
	if err != nil {
		return err
	}
	reverse(all)

	for _, node := range all {
		switch t {
		case typeTarget:
			target, err := g.MakeTarget(node, arch)
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "%s\n", target)
		case typeMakefileLine:
			entry, err := g.MakefileEntry(node)
			if err != nil {
				return err
			}
			if entry != "" {
				fmt.Fprintf(w, "%s\n", entry)
			}
		default:
			return fmt.Errorf("invalid type: %s", t)
		}
	}

	return nil
}
