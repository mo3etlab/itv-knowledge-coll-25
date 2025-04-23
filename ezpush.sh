
# echo $(git add -A)
# # echo $(git commit -am "update")
# echo $(git commit -m "update")
# echo $(git push --verbose)
# # echo $(git add -A && git commit -am "update" && git push --verbose)

git add -A


#origin_locate=$(pwd)

git add -A

current_time=$(date +'%Y-%m-%d %H:%M:%S')

commit_message="update when $current_time"
git commit -am "$commit_message"

##########
#git commit -am "update"

git push --verbose

# cd origin_locate
