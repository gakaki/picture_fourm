// MongoDB 数据库初始化脚本

// 连接到目标数据库
db = db.getSiblingDB('forum_db');

// 创建用户集合并设置索引
db.createCollection('users');
db.users.createIndex({ "username": 1 }, { unique: true });
db.users.createIndex({ "email": 1 }, { unique: true });
db.users.createIndex({ "created_at": -1 });

// 创建帖子集合并设置索引
db.createCollection('posts');
db.posts.createIndex({ "user_id": 1 });
db.posts.createIndex({ "created_at": -1 });
db.posts.createIndex({ "updated_at": -1 });
db.posts.createIndex({ "status": 1 });
db.posts.createIndex({ "tags": 1 });
db.posts.createIndex({ "like_count": -1 });
db.posts.createIndex({ "view_count": -1 });

// 创建评论集合并设置索引
db.createCollection('comments');
db.comments.createIndex({ "post_id": 1 });
db.comments.createIndex({ "user_id": 1 });
db.comments.createIndex({ "created_at": -1 });
db.comments.createIndex({ "parent_id": 1 });

// 创建图片生成记录集合并设置索引
db.createCollection('generations');
db.generations.createIndex({ "user_id": 1 });
db.generations.createIndex({ "created_at": -1 });
db.generations.createIndex({ "status": 1 });
db.generations.createIndex({ "prompt_hash": 1 });

// 创建模板集合并设置索引
db.createCollection('templates');
db.templates.createIndex({ "user_id": 1 });
db.templates.createIndex({ "created_at": -1 });
db.templates.createIndex({ "is_public": 1 });
db.templates.createIndex({ "category": 1 });
db.templates.createIndex({ "use_count": -1 });

// 创建用户积分交易记录集合并设置索引
db.createCollection('transactions');
db.transactions.createIndex({ "user_id": 1 });
db.transactions.createIndex({ "created_at": -1 });
db.transactions.createIndex({ "type": 1 });

// 创建点赞记录集合并设置索引
db.createCollection('likes');
db.likes.createIndex({ "user_id": 1, "post_id": 1 }, { unique: true });
db.likes.createIndex({ "post_id": 1 });
db.likes.createIndex({ "user_id": 1 });

// 创建关注记录集合并设置索引
db.createCollection('follows');
db.follows.createIndex({ "follower_id": 1, "following_id": 1 }, { unique: true });
db.follows.createIndex({ "follower_id": 1 });
db.follows.createIndex({ "following_id": 1 });

// 创建通知集合并设置索引
db.createCollection('notifications');
db.notifications.createIndex({ "user_id": 1 });
db.notifications.createIndex({ "created_at": -1 });
db.notifications.createIndex({ "is_read": 1 });
db.notifications.createIndex({ "type": 1 });

// 创建举报记录集合并设置索引
db.createCollection('reports');
db.reports.createIndex({ "reporter_id": 1 });
db.reports.createIndex({ "target_type": 1, "target_id": 1 });
db.reports.createIndex({ "created_at": -1 });
db.reports.createIndex({ "status": 1 });

print('数据库初始化完成！');
print('已创建的集合：');
db.runCommand("listCollections").cursor.firstBatch.forEach(
    function(collection) {
        print("- " + collection.name);
    }
);