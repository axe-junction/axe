import { useState, useEffect } from "react";
import {
  Heart,
  MessageCircle,
  Share2,
  Plus,
  MoreHorizontal,
} from "lucide-react";
import { useAppStore } from "../store/useAppStore";
import type { CommunityPost } from "../store/useAppStore";

const mockPosts: CommunityPost[] = [
  {
    id: "1",
    user: {
      id: "1",
      name: "Ahmed Benaissa",
      email: "ahmed@example.com",
      avatar:
        "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=150&h=150&fit=crop&crop=face",
    },
    content:
      "The new metro line from Bab Ezzouar to Alger Centre is amazing! Much faster than the old bus route. Highly recommended for daily commuters.",
    image:
      "https://images.unsplash.com/photo-1544620347-c4fd4a3d5957?w=400&h=300&fit=crop",
    likes: 24,
    comments: 8,
    timestamp: new Date("2024-01-15T10:30:00"),
    isLiked: false,
  },
  {
    id: "2",
    user: {
      id: "2",
      name: "Fatima Zerouali",
      email: "fatima@example.com",
      avatar:
        "https://images.unsplash.com/photo-1494790108755-2616b332c13c?w=150&h=150&fit=crop&crop=face",
    },
    content:
      "Pro tip: During rush hours, the tram is usually less crowded than buses. Plus, the views are better! 🚊",
    likes: 18,
    comments: 5,
    timestamp: new Date("2024-01-14T14:45:00"),
    isLiked: true,
  },
  {
    id: "3",
    user: {
      id: "3",
      name: "Yacine Bouhali",
      email: "yacine@example.com",
      avatar:
        "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=150&h=150&fit=crop&crop=face",
    },
    content:
      "Avoid taking the bus from Hydra to Centre during 5-6 PM. It takes forever due to traffic. Metro is much better.",
    likes: 31,
    comments: 12,
    timestamp: new Date("2024-01-13T16:20:00"),
    isLiked: false,
  },
];

export default function CommunityPage() {
  const { communityPosts, setCommunityPosts, toggleLike } = useAppStore();
  const [showNewPost, setShowNewPost] = useState(false);
  const [newPostContent, setNewPostContent] = useState("");

  useEffect(() => {
    if (communityPosts.length === 0) {
      setCommunityPosts(mockPosts);
    }
  }, [communityPosts, setCommunityPosts]);

  const handleLike = (postId: string) => {
    toggleLike(postId);
  };

  const handleShare = (post: CommunityPost) => {
    if (navigator.share) {
      navigator.share({
        title: "Transport Tip",
        text: post.content,
        url: window.location.href,
      });
    } else {
      // Fallback for browsers that don't support Web Share API
      navigator.clipboard.writeText(post.content);
      alert("Post copied to clipboard!");
    }
  };

  const handleNewPost = () => {
    if (newPostContent.trim()) {
      const newPost: CommunityPost = {
        id: Date.now().toString(),
        user: {
          id: "current-user",
          name: "You",
          email: "you@example.com",
          avatar:
            "https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=150&h=150&fit=crop&crop=face",
        },
        content: newPostContent,
        likes: 0,
        comments: 0,
        timestamp: new Date(),
        isLiked: false,
      };
      setCommunityPosts([newPost, ...communityPosts]);
      setNewPostContent("");
      setShowNewPost(false);
    }
  };

  const formatTimeAgo = (date: Date) => {
    const now = new Date();
    const diffInHours = Math.floor(
      (now.getTime() - date.getTime()) / (1000 * 60 * 60)
    );

    if (diffInHours < 1) return "Just now";
    if (diffInHours < 24) return `${diffInHours}h ago`;
    const diffInDays = Math.floor(diffInHours / 24);
    if (diffInDays < 7) return `${diffInDays}d ago`;
    const diffInWeeks = Math.floor(diffInDays / 7);
    return `${diffInWeeks}w ago`;
  };

  return (
    <div className="min-h-full bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow-sm border-b">
        <div className="p-4 flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Community</h1>
            <p className="text-gray-600">Share tips and experiences</p>
          </div>
          <button
            onClick={() => setShowNewPost(true)}
            className="bg-primary text-white p-3 rounded-full hover:bg-primary/90 transition-colors"
          >
            <Plus size={20} />
          </button>
        </div>
      </div>

      {/* New Post Modal */}
      {showNewPost && (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg w-full max-w-md">
            <div className="p-4 border-b">
              <h2 className="text-lg font-semibold">Share a tip</h2>
            </div>
            <div className="p-4">
              <textarea
                value={newPostContent}
                onChange={(e) => setNewPostContent(e.target.value)}
                placeholder="Share your transport experience or tip..."
                className="w-full h-32 p-3 border border-gray-300 rounded-lg resize-none focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              />
            </div>
            <div className="p-4 border-t flex gap-2">
              <button
                onClick={() => setShowNewPost(false)}
                className="flex-1 py-2 px-4 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                onClick={handleNewPost}
                className="flex-1 py-2 px-4 bg-primary text-white rounded-lg hover:bg-primary/90 transition-colors"
              >
                Post
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Posts */}
      <div className="max-w-2xl mx-auto p-4 space-y-4">
        {communityPosts.map((post) => (
          <div key={post.id} className="bg-white rounded-lg shadow-sm border">
            {/* Post Header */}
            <div className="p-4 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <img
                  src={post.user.avatar}
                  alt={post.user.name}
                  className="w-10 h-10 rounded-full object-cover"
                />
                <div>
                  <h3 className="font-medium text-gray-900">
                    {post.user.name}
                  </h3>
                  <p className="text-sm text-gray-500">
                    {formatTimeAgo(post.timestamp)}
                  </p>
                </div>
              </div>
              <button className="p-2 hover:bg-gray-100 rounded-full">
                <MoreHorizontal size={16} className="text-gray-400" />
              </button>
            </div>

            {/* Post Content */}
            <div className="px-4 pb-4">
              <p className="text-gray-900 mb-3">{post.content}</p>
              {post.image && (
                <img
                  src={post.image}
                  alt="Post"
                  className="w-full h-48 object-cover rounded-lg"
                />
              )}
            </div>

            {/* Post Actions */}
            <div className="px-4 py-3 border-t border-gray-100 flex items-center justify-between">
              <div className="flex items-center gap-4">
                <button
                  onClick={() => handleLike(post.id)}
                  className={`flex items-center gap-2 px-3 py-1 rounded-full transition-colors ${
                    post.isLiked
                      ? "text-red-500 bg-red-50"
                      : "text-gray-600 hover:text-red-500 hover:bg-red-50"
                  }`}
                >
                  <Heart
                    size={16}
                    fill={post.isLiked ? "currentColor" : "none"}
                  />
                  <span className="text-sm">{post.likes}</span>
                </button>
                <button className="flex items-center gap-2 px-3 py-1 rounded-full text-gray-600 hover:text-blue-500 hover:bg-blue-50 transition-colors">
                  <MessageCircle size={16} />
                  <span className="text-sm">{post.comments}</span>
                </button>
              </div>
              <button
                onClick={() => handleShare(post)}
                className="p-2 text-gray-600 hover:text-green-500 hover:bg-green-50 rounded-full transition-colors"
              >
                <Share2 size={16} />
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Empty State */}
      {communityPosts.length === 0 && (
        <div className="flex flex-col items-center justify-center py-12">
          <MessageCircle size={48} className="text-gray-400 mb-4" />
          <h3 className="text-lg font-medium text-gray-900 mb-2">
            No posts yet
          </h3>
          <p className="text-gray-600 text-center mb-4">
            Be the first to share a transport tip!
          </p>
          <button
            onClick={() => setShowNewPost(true)}
            className="bg-primary text-white px-6 py-2 rounded-lg hover:bg-primary/90 transition-colors"
          >
            Create Post
          </button>
        </div>
      )}
    </div>
  );
}
