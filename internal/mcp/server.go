package mcp

import (
	"context"
)

type Server struct {
	service ToolService
}

func NewServer(service ToolService) *Server {
	return &Server{service: service}
}

func (s *Server) HandleRequest(ctx context.Context, req Request) Response {
	switch req.Method {
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
	default:
		return Response{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &Error{Code: ErrCodeMethodNotFound, Message: "unknown tool: " + name},
		}
	}
}
