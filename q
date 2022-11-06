[33mcommit b4bffaeeb248b7e351437f2443616778da44deb3[m[33m ([m[1;36mHEAD -> [m[1;32minteract[m[33m, [m[1;31morigin/interact[m[33m)[m
Author: Shivam Tyagi <mergesuccessful@gmail.com>
Date:   Sun Nov 6 15:02:38 2022 +0530

    repo name fix

[1mdiff --git a/Jenkinsfile b/Jenkinsfile[m
[1mindex 6dc8291..0a4e5df 100644[m
[1m--- a/Jenkinsfile[m
[1m+++ b/Jenkinsfile[m
[36m@@ -5,7 +5,7 @@[m [mpipeline {[m
     stages {[m
         stage('Checkout Codebase'){[m
             steps{[m
[31m-                checkout scm: [$class: 'GitSCM', branches: [[name: '*/main']], userRemoteConfigs: [[credentialsId: 'PUBLIC_KEY', url: 'git@github.com:ShivamTyagi12345/Kritagya-Repository-Monitor-slackbot.git']]] [m
[32m+[m[32m                checkout scm: [$class: 'GitSCM', branches: [[name: '*/main']], userRemoteConfigs: [[credentialsId: 'PUBLIC_KEY', url: 'https://github.com/ShivamTyagi12345/Kritagya-Repository-Monitor-slackbot.git']]][m[41m [m
             }[m
         }[m
 [m
