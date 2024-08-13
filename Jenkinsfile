pipeline {
    agent any
    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                sh '/usr/local/go/bin/go mod tidy'
                sh '/usr/local/go/bin/go build ./cmd/bin/main.go'
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
                sshagent (credentials: ['jenkins-ssh-key']) {
                    sh 'ssh ${REMOTE_USER}@${REMOTE_SERVER} "systemctl stop ${SERVICE_CONSUMER_NAME}.service"'
                    sh 'ssh ${REMOTE_USER}@${REMOTE_SERVER} "systemctl stop ${SERVICE_NAME}.service"'
                    sh 'scp main ${REMOTE_USER}@${REMOTE_SERVER}:${FULLPATH_BINARY}'
                    sh 'ssh ${REMOTE_USER}@${REMOTE_SERVER} "systemctl start ${SERVICE_NAME}.service"'
                    sh 'ssh ${REMOTE_USER}@${REMOTE_SERVER} "systemctl start ${SERVICE_CONSUMER_NAME}.service"'
                }
            }
        }
    }
}