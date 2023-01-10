package main

import (
	"errors"
	"fmt"
	NDI "github.com/broadcastervc/ndigo"
	"go.uber.org/zap"
	"time"
)

type Router struct {
	log *zap.SugaredLogger

	inputs  map[string]*NDI.SourceType
	outputs map[string]*NDI.RoutingInstanceType
	matrix  map[string]string
}

func (r *Router) Init(logger *zap.SugaredLogger) error {
	r.log = logger.Named("Router")

	if !NDI.Initialize() {
		return errors.New("cannot init NDI")
	}

	r.inputs = make(map[string]*NDI.SourceType)
	r.outputs = make(map[string]*NDI.RoutingInstanceType, 0)
	r.matrix = make(map[string]string)

	return nil
}

type Source struct {
	Name string
	URL  string
}

func (r *Router) GetSources() []Source {
	log := r.log.Named("GetSources")

	log.Infow("getting sources")
	finder := NDI.FindCreateV2(NDI.NewFindCreateTypeRef(nil))

	_ = NDI.FindWaitForSources2(finder, 3000)
	time.Sleep(time.Second * 3)

	sources := NDI.FindGetCurrentSources2(finder)

	var list []Source
	for _, source := range sources {
		list = append(list, Source{
			Name: source.Name,
			URL:  source.URLAddress,
		})
	}

	return list
}

func (r *Router) InitInput(id string, name string) {
	log := r.log.Named("InitInput")
	log.Infow("init", "id", id, "ndiFullName", name)

	input := NDI.SourceType{
		PNdiName: name,
	}

	r.inputs[id] = &input
}

func (r *Router) InitOutput(id string, name string) {
	log := r.log.Named("InitOutput")
	log.Infow("init", "id", id, "ndiName", name)

	output := NDI.RoutingCreate(&NDI.RoutingCreateType{
		PNdiName: name,
		PGroups:  "public",
	})

	r.outputs[id] = &output
	r.matrix[id] = "" // TODO blank screen / placeholder image

	// cached assignment
	r.UpdateOutput(id, db.GetMatrix(id))
}

func (r *Router) UpdateOutput(outputID string, inputID string) error {
	log := r.log.Named("UpdateOutput").With("output", outputID, "input", inputID)
	log.Infow("updating")

	matrixID, found := r.matrix[outputID]
	if !found {
		return fmt.Errorf("output#%s not found in matrix", outputID)
	}
	if matrixID == inputID {
		if inputID != "" {
			log.Warnw("output is already specified input")
		}
		return nil
	}

	// find input/output
	output, found := r.outputs[outputID]
	if !found {
		return fmt.Errorf("output#%s not found in outputs", outputID)
	}
	input, found := r.inputs[inputID]
	if !found {
		return fmt.Errorf("input#%s not found in inputs", outputID)
	}

	// apply router
	ok := NDI.RoutingChange(*output, input)
	if !ok {
		return fmt.Errorf("NDI.RoutingChange not ok")
	}
	r.matrix[outputID] = inputID

	// cache assignment for restarts
	db.UpdateMatrix(outputID, inputID)

	return nil
}
