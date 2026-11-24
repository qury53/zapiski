package metadata_carrier

import "google.golang.org/grpc/metadata"

type GRPCCarrier struct {
	md *metadata.MD
}

func NewGRPCCarrier(md *metadata.MD) *GRPCCarrier {
	if md == nil {
		return &GRPCCarrier{md: &metadata.MD{}}
	}
	return &GRPCCarrier{md: md}
}

func (c *GRPCCarrier) Add(key, value string) {
	c.md.Append(key, value)
}

func (c *GRPCCarrier) Values(key string) []string {
	return c.md.Get(key)
}

func (c *GRPCCarrier) MD() metadata.MD {
	return c.md.Copy()
}
