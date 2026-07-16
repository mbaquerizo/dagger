package mcp

import (
	"context"

	"github.com/mbaquerizo/dagger/internal/issues"
	"github.com/mbaquerizo/dagger/internal/publish"
)

type Server struct {
	service ToolService
}

func NewServer(service ToolService) *Server {
	return &Server{service: service}
}

func (s *Server) HandleRequest(ctx context.Context, req Request) Response {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "notifications/initialized":
		return s.handleNotificationsInitialized()
	case "tools/list":
		return s.handleToolList(req)
	case "tools/call":
		return s.handleToolCall(ctx, req)
	default:
		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &Error{Code: ErrCodeMethodNotFound, Message: "method not found: " + req.Method},
		}
	}
}

func (s *Server) handleInitialize(req Request) Response {
	params, ok := req.Params.(map[string]interface{})

	if !ok {
		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &Error{Code: ErrCodeInvalidParams, Message: "invalid params"},
		}
	}

	clientProtocolVersion, _ := params["protocolVersion"].(string)

	switch clientProtocolVersion {
	case "2024-11-05":
	case "2025-03-26":
	case "2025-06-18":
	case "2025-11-25":
	default:
		clientProtocolVersion = "2025-11-25"
	}

	var tools struct{}

	result := InitializeResult{
		ProtocolVersion: clientProtocolVersion,
		ServerCapabilities: ServerCapabilities{
			Tools: tools,
		},
		ServerInfo: ServerInfo{
			Name:    "dagger-mcp",
			Version: "1.0.0",
		},
	}

	return Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}
}

func (s *Server) handleNotificationsInitialized() Response {
	return Response{}
}

func (s *Server) handleToolList(req Request) Response {
	return Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"tools": ListTools()},
	}
}

func (s *Server) handleToolCall(ctx context.Context, req Request) Response {
	params, ok := req.Params.(map[string]interface{})

	if !ok {
		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &Error{Code: ErrCodeInvalidParams, Message: "invalid params"},
		}
	}

	name, ok := params["name"].(string)

	if !ok || name == "" {
		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &Error{Code: ErrCodeInvalidParams, Message: "invalid or missing tool name"},
		}
	}

	args, ok := params["arguments"].(map[string]interface{})

	if !ok {
		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &Error{Code: ErrCodeInvalidParams, Message: "missing or invalid parameter: arguments"},
		}
	}

	switch name {
	case "get_issue":
		displayID, ok := args["display_id"].(string)

		if !ok || displayID == "" {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInvalidParams, Message: "missing or invalid parameter: display_id"},
			}
		}

		result, err := s.service.GetIssue(ctx, displayID)

		if err != nil {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInternal, Message: err.Error()},
			}
		}

		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  result,
		}
	case "get_doc":
		displayID, ok := args["display_id"].(string)

		if !ok || displayID == "" {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInvalidParams, Message: "missing or invalid parameter: display_id"},
			}
		}

		result, err := s.service.GetDoc(ctx, displayID)

		if err != nil {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInternal, Message: err.Error()},
			}
		}

		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  result,
		}
	case "list_issues":
		status, _ := args["status"].(string)

		result, err := s.service.ListIssues(ctx, status)

		if err != nil {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInternal, Message: err.Error()},
			}
		}

		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  result,
		}
	case "update_issue_status":
		displayID, ok := args["display_id"].(string)

		if !ok || displayID == "" {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInvalidParams, Message: "missing or invalid parameter: display_id"},
			}
		}

		newStatus, ok := args["status"].(string)

		if !ok || newStatus == "" {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInvalidParams, Message: "missing or invalid parameter: status"},
			}
		}

		result, err := s.service.UpdateIssueStatus(ctx, displayID, newStatus)

		if err != nil {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInternal, Message: err.Error()},
			}
		}

		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  result,
		}
	case "publish":
		publishType, ok := args["type"].(string)

		if !ok || publishType == "" {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInvalidParams, Message: "missing or invalid parameter: type"},
			}
		}

		title, ok := args["title"].(string)

		if !ok || title == "" {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInvalidParams, Message: "missing or invalid parameter: title"},
			}
		}

		body, ok := args["body"].(string)

		if !ok || body == "" {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInvalidParams, Message: "missing or invalid parameter: body"},
			}
		}

		projectID, ok := args["project_id"].(float64)

		if !ok || projectID == 0 {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInvalidParams, Message: "missing or invalid parameter: project_id"},
			}
		}

		publishReq := publish.PublishRequest{
			Type:      publishType,
			Title:     title,
			Body:      body,
			ProjectID: int(projectID),
		}

		parentIDRaw, ok := args["parent_id"].(float64)

		if ok && parentIDRaw != 0 {
			parentID := int(parentIDRaw)
			publishReq.ParentID = &parentID
		}

		metadataRaw, ok := args["metadata"].(map[string]interface{})

		if ok {
			var metadata publish.Metadata
			issueType, ok := metadataRaw["issue_type"].(string)

			if ok {
				metadata.IssueType = &issueType
			}

			status, ok := metadataRaw["status"].(string)

			if ok {
				metadata.Status = &status
			}

			tagsRaw, ok := metadataRaw["tags"].([]interface{})

			if ok {
				var tags []string
				for _, tagRaw := range tagsRaw {
					tag, ok := tagRaw.(string)

					if ok {
						tags = append(tags, tag)
					}
				}

				if len(tags) > 0 {
					metadata.Tags = tags
				}
			}

			relationshipsRaw, ok := metadataRaw["relationships"].([]interface{})

			if ok {
				var relationships []publish.Relationship

				for _, relationshipRaw := range relationshipsRaw {
					relationship, ok := relationshipRaw.(map[string]interface{})

					if ok {
						targetID, targetIDExists := relationship["target_id"].(float64)

						relationshipType, relationshipTypeExists := relationship["type"].(string)

						if !targetIDExists || !relationshipTypeExists {
							return Response{
								JSONRPC: "2.0",
								ID:      req.ID,
								Error:   &Error{Code: ErrCodeInvalidParams, Message: "invalid parameter: metadata.relationships"},
							}
						}

						relationships = append(relationships, publish.Relationship{
							TargetID: int(targetID),
							Type:     relationshipType,
						})
					}
				}

				metadata.Relationships = relationships
			}

			publishReq.Metadata = metadata
		}

		result, err := s.service.Publish(ctx, publishReq)

		if err != nil {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInternal, Message: err.Error()},
			}
		}

		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  result,
		}
	case "add_issue_relation":
		sourceID, ok := args["source_id"].(float64)

		if !ok || sourceID == 0 {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInvalidParams, Message: "missing or invalid parameter: source_id"},
			}
		}

		targetID, ok := args["target_id"].(float64)

		if !ok || targetID == 0 {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInvalidParams, Message: "missing or invalid parameter: target_id"},
			}
		}

		relationType, ok := args["relation_type"].(string)

		if !ok || relationType == "" {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInvalidParams, Message: "missing or invalid parameter: relation_type"},
			}
		}

		addIssueRelationReq := issues.AddIssueRelationRequest{
			SourceID:     int(sourceID),
			TargetID:     int(targetID),
			RelationType: relationType,
		}

		result, err := s.service.AddIssueRelation(ctx, addIssueRelationReq)

		if err != nil {
			return Response{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error:   &Error{Code: ErrCodeInternal, Message: err.Error()},
			}
		}

		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  result,
		}
	default:
		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &Error{Code: ErrCodeMethodNotFound, Message: "unknown tool: " + name},
		}
	}
}
