#!/bin/bash

# 确保脚本在任何错误发生时退出
set -e

# 获取第一次提交的哈希值
initial_commit=$(git rev-list --max-parents=0 HEAD)

# 获取最新一次提交的哈希值
latest_commit=$(git rev-parse HEAD)

# 创建一个新的分支来操作，避免直接在主分支上进行危险操作
git checkout -b temp-branch

# 重置到最新一次提交
git reset --hard $latest_commit

# 提交当前最新代码
git commit --allow-empty -m "Keep latest commit"

# 强制更新主分支只保留第一次和最新一次提交
git checkout main
git reset --hard $initial_commit
git cherry-pick $latest_commit

# 强制推送到远程仓库
git push origin main --force

# 删除临时分支
git branch -d temp-branch

echo "Successfully retained only the initial and latest commit."
