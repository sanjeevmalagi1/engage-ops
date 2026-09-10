package customerio

import (
	"context"
	"fmt"

	"github.com/sanjeevmalagi1/engage-ops/internal/core"
	"github.com/sanjeevmalagi1/engage-ops/internal/core/diff"
)

// segmentResource manages Customer.io manual segments
// (https://customer.io/docs/api/app/#tag/segments).
type segmentResource struct {
	provider *Provider
}

type segmentAPIModel struct {
	Segment struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"segment"`
}

func (r *segmentResource) Diff(_ context.Context, addr core.Address, state, desired core.Attributes) (*core.Change, error) {
	return diff.Compute(addr, state, desired), nil
}

func (r *segmentResource) Apply(ctx context.Context, change *core.Change) (core.Attributes, error) {
	switch change.Action {
	case core.ActionCreate:
		return r.create(ctx, change.After)
	case core.ActionUpdate:
		return r.update(ctx, change.Before, change.After)
	case core.ActionDelete:
		return nil, r.delete(ctx, change.Before)
	case core.ActionNoop:
		return change.Before, nil
	default:
		return nil, fmt.Errorf("customerio segment: unsupported action %q", change.Action)
	}
}

func (r *segmentResource) Read(ctx context.Context, state core.Attributes) (core.Attributes, error) {
	id, ok := state["_id"]
	if !ok {
		return nil, fmt.Errorf("customerio segment: state missing _id")
	}

	var resp segmentAPIModel
	if err := r.provider.request(ctx, "GET", fmt.Sprintf("/segments/%v", id), nil, &resp); err != nil {
		return nil, nil // treat any read failure as "no longer exists" for drift detection
	}

	return core.Attributes{
		"_id":  resp.Segment.ID,
		"name": resp.Segment.Name,
	}, nil
}

func (r *segmentResource) create(ctx context.Context, desired core.Attributes) (core.Attributes, error) {
	name, _ := desired["name"].(string)
	if name == "" {
		return nil, fmt.Errorf("customerio segment: %q is required", "name")
	}

	payload := map[string]any{"segment": map[string]any{"name": name}}
	var resp segmentAPIModel
	if err := r.provider.request(ctx, "POST", "/segments", payload, &resp); err != nil {
		return nil, err
	}

	return core.Attributes{"_id": resp.Segment.ID, "name": name}, nil
}

func (r *segmentResource) update(ctx context.Context, state, desired core.Attributes) (core.Attributes, error) {
	id, ok := state["_id"]
	if !ok {
		return nil, fmt.Errorf("customerio segment: state missing _id")
	}
	name, _ := desired["name"].(string)

	payload := map[string]any{"segment": map[string]any{"name": name}}
	if err := r.provider.request(ctx, "PUT", fmt.Sprintf("/segments/%v", id), payload, nil); err != nil {
		return nil, err
	}

	return core.Attributes{"_id": id, "name": name}, nil
}

func (r *segmentResource) delete(ctx context.Context, state core.Attributes) error {
	id, ok := state["_id"]
	if !ok {
		return fmt.Errorf("customerio segment: state missing _id")
	}
	return r.provider.request(ctx, "DELETE", fmt.Sprintf("/segments/%v", id), nil, nil)
}
