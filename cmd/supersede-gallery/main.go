// supersede-gallery is a command line tool to clone and supersede an existing gallery record assigning
// updated parent and hierarchy information at the same time. For example:
// 	$> ./bin/supersede-gallery -architecture-reader-uri repo:///usr/local/build/collection/sfomuseum-data-architecture/ -gallery-id 1763595133 -parent-id 1763588365
package main

import (
	"context"
	"flag"
	sfom_reader "github.com/sfomuseum/go-sfomuseum-reader"
	sfom_writer "github.com/sfomuseum/go-sfomuseum-writer"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-export/v2"
	"github.com/whosonfirst/go-whosonfirst-id"	
	"github.com/whosonfirst/go-writer"
	"log"
	
)

func main() {

	architecture_reader_uri := flag.String("architecture-reader-uri", "repo:///usr/local/data/sfomuseum-data-architecture", "")
	architecture_writer_uri := flag.String("architecture-writer-uri", "", "If empty, the value of the -architecture-reader-uri flag will be used.")

	gallery_id := flag.Int64("gallery-id", 0, "The SFO Museum gallery ID to supersede")
	parent_id := flag.Int64("parent-id", 0, "The SFO Museum parent ID of the new gallery")	

	name := flag.String("name", "", "An optional name for the new gallery. If empty the name of the previous gallery will be used.")
	map_id := flag.String("map_id", "", "An optional map ID for the new gallery. If empty the name of the previous gallery will be used.")

	flag.Parse()
	
	ctx := context.Background()

	if *architecture_writer_uri == "" {
		*architecture_writer_uri = *architecture_reader_uri
	}

	arch_r, err := reader.NewReader(ctx, *architecture_reader_uri)

	if err != nil {
		log.Fatalf("Failed to create architecture reader, %v", err)
	}

	arch_wr, err := writer.NewWriter(ctx, *architecture_writer_uri)

	if err != nil {
		log.Fatalf("Failed to create architecture writer, %v", err)
	}

	id_provider, err := id.NewProvider(ctx)

	if err != nil {
		log.Fatalf("Failed to create ID provider, %v", err)
	}
		
	gallery_f, err := sfom_reader.LoadBytesFromID(ctx, arch_r, *gallery_id)

	if err != nil {
		log.Fatalf("Failed to load gallery record, %v", err)
	}

	parent_f, err := sfom_reader.LoadBytesFromID(ctx, arch_r, *parent_id)

	if err != nil {
		log.Fatalf("Failed to load parent record, %v", err)
	}

	new_id, err := id_provider.NewID()

	if err != nil {
		log.Fatalf("Failed to create new ID, %v", err)
	}

	new_updates := map[string]interface{}{
		"properties.id": new_id,		
		"properties.wof:id": new_id,
		"properties.wof:parent_id": *parent_id,
		"properties.wof:hierarchy": gjson.GetBytes(parent_f, "properties.wof:hierarchy").Value(),
		"properties.mz:is_current": gjson.GetBytes(parent_f, "properties.mz:is_current").Value(),
		"properties.edtf:inception": gjson.GetBytes(parent_f, "properties.edtf:inception").Value(),
		"properties.edtf:cessation": gjson.GetBytes(parent_f, "properties.edtf:cessation").Value(),
		"properties.wof:supersedes": []int64{ *gallery_id },
	}

	if *name != "" {
		new_updates["properties.wof:name"] = *name
	}

	if *map_id != "" {
		new_updates["properties.sfomuseum:map_id"] = *map_id
	}

	// Create and record the new gallery
	
	_, new_gallery, err := export.AssignPropertiesIfChanged(ctx, gallery_f, new_updates)

	if err != nil {
		log.Fatalf("Failed to export new gallery, %v", err)
	}

	_, err = sfom_writer.WriteFeatureBytes(ctx, arch_wr, new_gallery)

	if err != nil {
		log.Fatalf("Failed to write new gallery, %v", err)
	}

	old_updates := map[string]interface{}{
		"properties.wof:superseded_by": []int64{ new_id },
	}

	// Now update the previous gallery
	
	_, gallery_f, err = export.AssignPropertiesIfChanged(ctx, gallery_f, old_updates)

	if err != nil {
		log.Fatalf("Failed to export new gallery, %v", err)
	}

	_, err = sfom_writer.WriteFeatureBytes(ctx, arch_wr, gallery_f)

	if err != nil {
		log.Fatalf("Failed to write previous gallery, %v", err)
	}

	log.Printf("Created new gallery record with ID %d\n", new_id)
}
