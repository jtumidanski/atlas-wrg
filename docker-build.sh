if [ $1 = NO-CACHE ]
then
   docker build --no-cache --tag atlas-wrg:latest .
else
   docker build --tag atlas-wrg:latest .
fi
