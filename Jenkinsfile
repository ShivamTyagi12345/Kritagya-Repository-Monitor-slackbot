pipeline {

    agent any 

    stages {
        stage('Checkout Codebase'){
            steps{
                checkout scm: [$class: 'GitSCM', branches: [[name: '*/main']], userRemoteConfigs: [[credentialsId: 'github-ssh-key', url: 'git@github.com:ShivamTyagi12345/Kritagya-Repository-Monitor-slackbot.git']]] 
            }
        }

        stage('Build') {
            steps {
                echo 'Building Codebase'
            }                           
        }

        stage('Test'){
            steps {
                echo 'Running Tests on changes'
            }
        }
        
        stage('Deploy'){
            steps{
                echo 'Done!'
            }
        }
    }

    post {

        always {
            echo 'Sending Slack message'
            sh "go run main.go ${BUILD_URL} ${currentBuild.currentResult} ${BUILD_NUMBER} ${JOB_NAME} "
        }
    }
}