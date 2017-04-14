# maybe more powerful
# for mac (sed for linux is different)
dir=`echo ${PWD##*/}`
grep "robot-rmdmy" * -R | grep -v Godeps | awk -F: '{print $1}' | sort | uniq | xargs sed -i '' "s#robot-rmdmy#$dir#g"
mv robot-rmdmy.ini $dir.ini

