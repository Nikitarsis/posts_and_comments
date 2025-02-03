package main

import (
	"github.com/Nikitarsis/posts_and_comments/messages"
	mtd "github.com/Nikitarsis/posts_and_comments/muted_posts"
	onlyhttp "github.com/Nikitarsis/posts_and_comments/only_http"
	cap "github.com/Nikitarsis/posts_and_comments/posts_with_comms"
	tdao "github.com/Nikitarsis/posts_and_comments/translationdao"
	"github.com/Nikitarsis/posts_and_comments/users"
)

var posts cap.PostHypervisor
var authors users.AuthorManager
var mutedPosts mtd.MutedPost
var contentManager messages.MessagesController

func daoFromIPostWithMessage(post cap.IPost, userId users.IUser, message *string) tdao.PostDao {
	return tdao.PostDao{
		PostId:  post.GetMessageId().GetId(),
		UserId:  userId.GetId(),
		Message: message,
	}
}

func daoComment(comment cap.IPost, userId users.IUser) tdao.CommentDao {
	post := comment.GetMessageId()
	parent, _ := comment.GetParentId()
	children := comment.GetChildrenIds()
	retChildren := make([]uint64, len(children))
	for _, child := range children {
		retChildren = append(retChildren, child.GetId())
	}
	return tdao.CommentDao{
		UserId:      userId.GetId(),
		PostId:      post.GetId(),
		ParentId:    parent.GetId(),
		ChildrenIds: retChildren,
		Message:     nil,
	}
}

func daoCommentWithMessage(comment cap.IPost, userId users.IUser, message *string) tdao.CommentDao {
	post := comment.GetMessageId()
	parent, _ := comment.GetParentId()
	children := comment.GetChildrenIds()
	retChildren := make([]uint64, len(children))
	for _, child := range children {
		retChildren = append(retChildren, child.GetId())
	}
	return tdao.CommentDao{
		UserId:      userId.GetId(),
		PostId:      post.GetId(),
		ParentId:    parent.GetId(),
		ChildrenIds: retChildren,
		Message:     message,
	}
}

func listPosts() []tdao.PostDao {
	postSlice := posts.GetPosts()
	dao := make([]tdao.PostDao, len(postSlice))
	for i, pst := range postSlice {
		author, check := authors.GetAuthorOfPost(pst)
		if !check {
			author = users.GetNullUsr()
		}
		dao[i] = tdao.PostDao{
			PostId:  pst.GetId(),
			UserId:  author.GetId(),
			Message: nil,
		}
	}
	return dao
}

func mutePost(postId uint64, userId uint64) tdao.PROBLEM {
	post := messages.GetMessageId(postId)
	user := users.GetUser(userId)
	if !authors.CheckAuthor(user) {
		return tdao.NO_SUCH_USER
	}
	author, check := authors.GetAuthorOfPost(post)
	if !check {
		return tdao.INCORRECT_USER
	}
	if userId != author.GetId() {
		return tdao.INCORRECT_USER
	}
	mutedPosts.ForbidComment(post)
	return tdao.NO_PROBLEM
}

func unmutePost(postId uint64, userId uint64) tdao.PROBLEM {
	post := messages.GetMessageId(postId)
	user := users.GetUser(userId)
	if !authors.CheckAuthor(user) {
		return tdao.NO_SUCH_USER
	}
	author, check := authors.GetAuthorOfPost(post)
	if !check {
		return tdao.INCORRECT_USER
	}
	if userId != author.GetId() {
		return tdao.INCORRECT_USER
	}
	mutedPosts.AllowComment(post)
	return tdao.NO_PROBLEM
}

func getPost(post uint64, fromId uint64, toId uint64) (tdao.PostDao, []tdao.CommentDao, tdao.PROBLEM) {
	postId := messages.GetMessageId(post)
	str := contentManager.GetContent(postId)
	comPost, check := posts.GetPost(postId)
	if !check {
		return tdao.GetEmptyPostDao(), make([]tdao.CommentDao, 0), tdao.NO_SUCH_POST
	}
	comments, err := comPost.GetCommentPage(int(fromId), int(toId))
	if err != nil {
		return tdao.GetEmptyPostDao(), make([]tdao.CommentDao, 0), tdao.NO_SUCH_POST
	}
	retCom := make([]tdao.CommentDao, len(comments))
	user, exists := authors.GetAuthorOfPost(comPost.GetPost().GetMessageId())
	if !exists {
		return tdao.GetEmptyPostDao(), make([]tdao.CommentDao, 0), tdao.NO_SUCH_POST
	}
	retPost := daoFromIPostWithMessage(comPost.GetPost(), user, str)
	for i, comment := range comments {
		user, exists := authors.GetAuthorOfPost(comment.GetMessageId())
		if !exists {
			return tdao.GetEmptyPostDao(), make([]tdao.CommentDao, 0), tdao.NO_SUCH_USER
		}
		retCom[i] = daoComment(comment, user)
	}
	return retPost, retCom, tdao.NO_PROBLEM
}

func createPost(user uint64, message *string) (uint64, tdao.PROBLEM) {
	postId := messages.GetNewMessageId()
	posts.NewPost(postId)
	return postId.GetId(), tdao.NO_PROBLEM
}

func updatePost(post uint64, user uint64, message *string) tdao.PROBLEM {
	postId := messages.GetMessageId(post)
	userId := users.GetUser(user)
	if !posts.HasPost(postId) {
		return tdao.NO_SUCH_POST
	}
	author, _ := authors.GetAuthorOfPost(postId)
	if author != userId {
		return tdao.NO_SUCH_USER
	}
	contentManager.SetContent(postId, message)
	return tdao.NO_PROBLEM
}

func deletePost(post uint64, user uint64) tdao.PROBLEM {
	userId := users.GetUser(user)
	postId := messages.GetMessageId(post)
	author, check := authors.GetAuthorOfPost(postId)
	if !check {
		return tdao.INCORRECT_USER
	}
	if author != userId {
		return tdao.INCORRECT_USER
	}

	posts.DeletePost(postId)
	contentManager.DeleteContent(postId)
	return tdao.NO_PROBLEM
}

func getComment(comment uint64) (tdao.PROBLEM, *tdao.CommentDao) {
	commentId := messages.GetMessageId(comment)
	userId, hasAuthor := authors.GetAuthorOfPost(commentId)
	if !hasAuthor {
		userId = users.GetNullUsr()
	}
	com, check := posts.GetComment(commentId)
	if !check {
		return tdao.NO_SUCH_POST, nil
	}
	message := contentManager.GetContent(commentId)
	ret := daoCommentWithMessage(com, userId, message)
	return tdao.NO_PROBLEM, &ret
}

func createComment(user uint64, parent uint64, message *string) (uint64, tdao.PROBLEM) {
	postId := messages.GetNewMessageId()
	parentId := messages.GetMessageId(parent)
	posts.NewComment(parentId, postId)
	contentManager.SetContent(postId, message)
	return postId.GetId(), tdao.NO_PROBLEM
}

func updateComment(comment uint64, user uint64, message *string) tdao.PROBLEM {
	commentId := messages.GetMessageId(comment)
	userId := users.GetUser(user)
	if !posts.HasPost(commentId) {
		return tdao.NO_SUCH_POST
	}
	author, _ := authors.GetAuthorOfPost(commentId)
	if author != userId {
		return tdao.NO_SUCH_USER
	}
	contentManager.SetContent(commentId, message)
	return tdao.NO_PROBLEM
}

func deleteComment(comment uint64, user uint64) tdao.PROBLEM {
	userId := users.GetUser(user)
	commentId := messages.GetMessageId(comment)
	author, check := authors.GetAuthorOfPost(commentId)
	if !check {
		return tdao.INCORRECT_USER
	}
	if author != userId {
		return tdao.INCORRECT_USER
	}

	contentManager.DeleteContent(commentId)
	return tdao.NO_PROBLEM
}

func main() {
	posts = cap.NewPostHypervisor()
	authors = users.NewAuthorManager()
	mutedPosts = mtd.NewMutedPost()
	contentManager = messages.NewMesssagesController()

	postsCallbacks := onlyhttp.PostsCallback{
		ListPosts: listPosts,
	}
	muteCallbacks := onlyhttp.MutePostCallback{
		MutePost: mutePost,
	}
	unmuteCallbacks := onlyhttp.UnmutePostCallback{
		UnmutePost: unmutePost,
	}
	postCallbacks := onlyhttp.PostCallback{
		GetPost:    getPost,
		CreatePost: createPost,
		UpdatePost: updatePost,
		DeletePost: deletePost,
	}
	commentCallbacks := onlyhttp.CommentCallback{
		GetComment:    getComment,
		CreateComment: createComment,
		UpdateComment: updateComment,
		DeleteComment: deleteComment,
	}

	serverCalls := onlyhttp.ServerCallbacks{
		Log:          func(s string) {},
		Posts_:       postsCallbacks,
		Post_:        postCallbacks,
		Post_mute_:   muteCallbacks,
		Post_unmute_: unmuteCallbacks,
		Comment_:     commentCallbacks,
	}
	onlyhttp.StartServer(serverCalls)
}
