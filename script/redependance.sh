cd /
tar zcf dep.tar.gz /projectdep
rm -rf /tmp/dep
mkdir -p /tmp/dep
mv /dep.tar.gz /tmp/dep/
cd /tmp/dep
git init
git remote add origin https://git.oschina.net/lijianying10/mygodependance.git
git add --all
git commit -m 'update'
git push origin master --force
