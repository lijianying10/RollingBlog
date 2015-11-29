cd philotop
rm -rf *
cp -rf ../public/* .
git add --all
git commit -m 'update'
git push -f origin master 
