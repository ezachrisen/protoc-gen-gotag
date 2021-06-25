package module

const firestoreTpl = `
// Code generated by protoc-{{name}}. DO NOT EDIT.
// versions:
// source: {{ .InputPath }} 

package {{ package . }} 

import (
 "fmt"
 "context"
 "cloud.google.com/go/firestore"
 "github.com/golang/protobuf/ptypes"
 "google.golang.org/genproto/protobuf/field_mask"
  "google.golang.org/api/iterator"
  "alticeusa.com/maui/metadata"
)

{{ range .Messages }}
  {{if shouldGenerateFirestore . }}

   func Get{{.Name}}(ctx context.Context, client *firestore.Client, collection, name string) (*{{.Name}}, error) {

	   if client == nil {
		   return nil, fmt.Errorf("firestore client is nil")
	   }

	   d, err := client.Collection(collection).Doc(name).Get(ctx)
	   if err != nil {
		   return nil, err
	   }

	   var doc {{.Name}}
	   if err := d.DataTo(&doc); err != nil {
		   return nil, fmt.Errorf("error converting Firestore document to %T: %v", doc, err)
	   }

	   doc.Name = name
	   return &doc, nil
   }


 func Create{{.Name}}(ctx context.Context, client *firestore.Client, collection string, doc *{{.Name}}) error {

	   if client == nil {
		   return fmt.Errorf("Create{{.Name}}: firestore client is nil")
	   }

	doc.UpdatedBy = metadata.GetUsernameFromMetadata(ctx)
	doc.Created = ptypes.TimestampNow()

	_, err := client.Collection(collection).Doc(doc.Name).Set(ctx, doc)
	if err != nil {
      return  fmt.Errorf("Create{{.Name}}: %w", err)
	}

	return  nil
 }

 func Delete{{.Name}}(ctx context.Context, client *firestore.Client, collection, name string) error {
	   if client == nil {
		   return fmt.Errorf("firestore client is nil")
	   }

	_, err := client.Collection(collection).Doc(name).Delete(ctx)
	if err != nil {
      return fmt.Errorf("Delete{{.Name}}: %w", err)
	}
    return nil 
}


 func Update{{.Name}}(ctx context.Context, client *firestore.Client, collection string, req *{{.Name}}UpdateRequest) error {
	   if client == nil {
		   return fmt.Errorf("firestore client is nil")
	   }

	// Add the system-mandated fields to the field mask
	mask := field_mask.FieldMask{
		Paths: []string{"updated_by", "updated"},
	}
    
    if req.UpdateMask == nil {
       return fmt.Errorf("Update{{.Name}}: field mask is nil")
   }

    // TODO: prevent unauthorized update of the UpdatedBy and Updated field
	// Only allow some accounts the ability to override the system-generated fields,
	/// such as for batch loading / transfers, etc. For now, block any outside updates.

	for _, p := range req.UpdateMask.Paths {
		for _, x := range  []string{"updated_by", "updated", "name"} {
			if p == x {
				return fmt.Errorf("Update{{.Name}} updating '%s': Updates to '%s' not permitted", req.{{.Name}}.Name, x)
			}
		}
	}

    // TODO: Proto field masks can be "a.b", Firestore expects ["a","b"]. This difference doesn't matter for rules.
    mask.Paths = append(mask.Paths, req.UpdateMask.Paths...)

	fps := []firestore.FieldPath{}
	for _, p := range mask.Paths {
		fps = append(fps, []string{p})
	}

	req.Rule.UpdatedBy = metadata.GetUsernameFromMetadata(ctx)
	req.Rule.Updated = ptypes.TimestampNow()

	_, err := client.Collection(collection).Doc(req.{{.Name}}.Name).Set(ctx, req.{{.Name}}, firestore.Merge(fps...))

	if err != nil {
		return fmt.Errorf("Update{{.Name}}: error saving '%s': %v", req.{{.Name}}.Name, err)
	} 
    return nil 
  }


func List{{.Name}}s(ctx context.Context, client *firestore.Client, collection string, req *{{.Name}}ListRequest) (*{{.Name}}List, error) {

	   if client == nil {
		   return nil, fmt.Errorf("firestore client is nil")
	   }

	limit := 50

	if limit > 50 || limit < 0 {
		limit = 50
	}

	pl := {{.Name}}List{}

	pl.{{.Name}}s = make([]*{{.Name}}, 0, limit)

	query := client.Collection(collection)

	iter := query.Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var p {{.Name}}
		if err := doc.DataTo(&p); err != nil {
             return nil, fmt.Errorf("converting Firestore document %T: %v", p, err)
		}
		p.Name = doc.Ref.ID
		pl.{{.Name}}s = append(pl.{{.Name}}s, &p)
	}
	return &pl, nil
  }



 {{ end }} 
{{ end }}
`

/*



 */
