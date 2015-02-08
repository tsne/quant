// Package quant provides the tools to measure application metrics.
//
// The first step to application metrics is to create a registry.
// A registry provides functions to create several metric types.
// Each metric has a unique name within the registry which is
// used to identify the metric for retrieval.
//
// After creating the registry and using the metrics a snapshot
// of each metric can be written to specified locations. This could
// be achieved with the Report function of Registry which takes a
// collection of reporters. Each reporter processes all the snapshots
// the registry provides. By calling Report regularly a metrics stream
// can be created to constantly measure application characteristics.
// This is where the Reporting comes into play. With a reporting
// the Report function is called periodically at a specified interval
// and to the specified reporters.
package quant
