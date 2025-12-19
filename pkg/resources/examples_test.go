// [REQ:RESOURCE_MANAGEMENT] Resource management requirement
// [ARCH:RESOURCE_MANAGEMENT] Resource manager architecture
// [IMPL:RESOURCE_MANAGER] Resource manager implementation
package resources

import (
	"fmt"
)

// nopResource is a simple example implementation of the Resource interface
// used by the godoc examples in this package.
type nopResource struct {
	name string
}

func (n *nopResource) Cleanup() error { return nil }
func (n *nopResource) String() string { return "nop:" + n.name }

// ExampleResourceManager_basic demonstrates basic registration and inspection
// of resources without performing filesystem operations. It uses AddTempFile
// and AddTempDir to register resources and then inspects the count.
func ExampleResourceManager_basic() {
	rm := NewResourceManager()
	// Register resources (note: these paths are examples and no filesystem
	// operations are performed here; Cleanup is not called in this example)
	rm.AddTempFile("/tmp/example-file.txt")
	rm.AddTempDir("/tmp/example-dir")

	fmt.Println(rm.GetResourceCount())
	// Output: 2
}

// ExampleResourceManager_CustomResource shows how to implement the Resource
// interface for a custom type and register it with the ResourceManager.
func Example_customresource() {
	rm := NewResourceManager()
	rm.AddResource(&nopResource{name: "demo"})

	for _, r := range rm.GetResources() {
		fmt.Println(r.String())
	}
	// Output: nop:demo
}
