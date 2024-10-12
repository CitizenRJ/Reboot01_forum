package like

import (
	"database/sql"
	"encoding/json"
	"errors"
	e "forum/error"
	"log/slog"
	"net/http"
)

type LikesController struct {
	s LikesService
}

func NewLikesController(s LikesService) *LikesController {
	return &LikesController{s: s}
}

func (c *LikesController) LikeDislikePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.HandleMethod(w, r)
		return
	}

	var req InteractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(r.Context(), "Error decoding request body", "err", err)
		e.HandleInternalError(w, r)
		return
	}

	if req.PostID == nil {
		e.HandleBadRequest(w, r)
		return
	}

	token, err := r.Cookie("session_id")
	if err != nil {
		slog.ErrorContext(r.Context(), "Error getting the token", "err", err)
		e.HandleInternalError(w, r)
		return
	}

	userID, err := c.s.GetUserIDFromSession(r.Context(), token.Value)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error getting user_id from session", "err", err)
		e.HandleInternalError(w, r)
		return
	}

	like, err := c.s.CheckPostInteractions(r.Context(), userID, *req.PostID)
	if errors.Is(err, sql.ErrNoRows) {
		if err = c.s.InteractWithPost(r.Context(), userID, *req.PostID, req.IsLike); err != nil {
			slog.ErrorContext(r.Context(), "Error in interacting with post", "err", err)
			e.HandleInternalError(w, r)
			return
		}
	} else if like.IsLike == req.IsLike {
		if err = c.s.RemovePostInteraction(r.Context(), userID, *req.PostID); err != nil {
			slog.ErrorContext(r.Context(), "Error removing post interaction", "err", err)
			e.HandleInternalError(w, r)
			return
		}
	} else {
		if err = c.s.RemovePostInteraction(r.Context(), userID, *req.PostID); err != nil {
			slog.ErrorContext(r.Context(), "Error removing post interaction", "err", err)
			e.HandleInternalError(w, r)
			return
		}
		if err = c.s.InteractWithPost(r.Context(), userID, *req.PostID, req.IsLike); err != nil {
			slog.ErrorContext(r.Context(), "Error creating a new interaction with post", "err", err)
			e.HandleInternalError(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Interaction done successfully"}); err != nil {
		slog.ErrorContext(r.Context(), "Error decoding request body", "err", err)
		e.HandleInternalError(w, r)
	}
}

func (c *LikesController) InteractWithComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		e.HandleMethod(w, r)
		return
	}

	var req InteractRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		e.HandleBadRequest(w, r)
		return
	}

	if req.CommentID == nil {
		e.HandleBadRequest(w, r)
		return
	}

	token, err := r.Cookie("session_id")
	if err != nil {
		slog.ErrorContext(r.Context(), "Error getting the token", "err", err)
		e.HandleInternalError(w, r)
		return
	}

	userID, err := c.s.GetUserIDFromSession(r.Context(), token.Value)
	if err != nil {
		slog.ErrorContext(r.Context(), "Error getting user_id from session", "err", err)
		e.HandleInternalError(w, r)
		return
	}

	like, err := c.s.CheckCommentInteractions(r.Context(), userID, *req.CommentID)
	if errors.Is(err, sql.ErrNoRows) {
		if err = c.s.InteractWithComment(r.Context(), userID, *req.CommentID, req.IsLike); err != nil {
			slog.ErrorContext(r.Context(), "Error in interacting with comment", "err", err)
			e.HandleInternalError(w, r)
			return
		}
	} else if like.IsLike == req.IsLike {
		if err = c.s.RemoveCommentInteraction(r.Context(), userID, *req.CommentID); err != nil {
			slog.ErrorContext(r.Context(), "Error removing comment interaction", "err", err)
			e.HandleInternalError(w, r)
			return
		}
	} else {
		if err = c.s.RemoveCommentInteraction(r.Context(), userID, *req.CommentID); err != nil {
			slog.ErrorContext(r.Context(), "Error removing comment interaction", "err", err)
			e.HandleInternalError(w, r)
			return
		}
		if err = c.s.InteractWithComment(r.Context(), userID, *req.CommentID, req.IsLike); err != nil {
			slog.ErrorContext(r.Context(), "Error creating a new interaction with comment", "err", err)
			e.HandleInternalError(w, r)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Interaction done successfully"}); err != nil {
		slog.ErrorContext(r.Context(), "Error decoding request body", "err", err)
		e.HandleInternalError(w, r)
	}
}

func (c *LikesController) GetInteractions(w http.ResponseWriter, r *http.Request) {
	var req GetInteractionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.ErrorContext(r.Context(), "Error decoding request body", "err", err)
		e.HandleBadRequest(w, r)
		return
	}

	var resp GetInteractionsResponse
	var err error
	if req.PostID != nil {
		resp, err = c.s.GetPostsInteractions(r.Context(), *req.PostID)
		if err != nil {
			e.HandleInternalError(w, r)
			return
		}
	} else if req.CommentID != nil {
		resp, err = c.s.GetCommentsInteractions(r.Context(), *req.CommentID)
		if err != nil {
			e.HandleInternalError(w, r)
			return
		}
	} else {
		e.HandleBadRequest(w, r)
		return
	}

	response := map[string]int{"likes": resp.Likes, "dislikes": resp.Dislikes}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.ErrorContext(r.Context(), "Error encoding json", "err", err)
		e.HandleInternalError(w, r)
	}
}
