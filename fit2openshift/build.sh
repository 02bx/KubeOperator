#! /bin/bash


red=31
green=32
yellow=33
blue=34

function colorMsg()
{
  echo -e "\033[$1m $2 \033[0m"
}
logPath="/opt/fit2openshift/logs/install/"
timestamp=$(date -d now +%F)
errorLogFile=${logPath}"error/install_error_"${timestamp}".log"
infoLogFile=${logPath}"info/install_info_"${timestamp}".log"
fullLogFile=${logPath}"install_"${timestamp}".log"
printf "%-65s .......... " "Build fit2openshift webconsole ui:"
cd ui
npm install >/dev/null 2>&1
ng build --prod 1>>$infoLogFile 2>>$errorLogFile
if [ "$?" == "0" ];then
    colorMsg $green "[OK]"
else
    colorMsg $red "[DEFEATED]"
    printf "\n"
    printf "Build fit2openshift webconsole ui  defeated! An error log in :"${errorLogFile}
    printf "\n"
    exit 1
fi
printf "\n"
printf "%-65s .......... " "Build fit2openshift webconsole api: "
cd .. && docker build --rm=true --tag=registry.fit2cloud.com/fit2anything/fit2openshift/fit2openshift-app:latest . 1>>$infoLogFile 2>>$errorLogFile

if [ "$?" == "0" ];then
    colorMsg $green "[OK]"
else
    colorMsg $red "[DEFEATED]"
    printf "\n"
    printf "Build fit2openshift webconsole api  defeated! An error log in :"${errorLogFile}
    printf "\n"
    exit 1
fi

