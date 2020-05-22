pipeline{
    agent any
    tools {
        go 'go-1.14.3'
    }
    environment {
        GO111MODULE = 'on'
    }
    stages{
        stage("Clone"){
            steps{
                echo "Cloning the git repository"
                git branch: 'master', url: 'https://github.com/vocacorg/terraform-provider-github.git'
                echo "Content in working directory"
                sh "ls -la ."
            }
            post{
                always{
                    echo "========always========"
                }
                success{
                    echo "========A executed successfully========"
                }
                failure{
                    echo "========A execution failed========"
                }
            }
        }
        stage("Code Quality"){
            environment {
                scannerHome = tool 'SonarqubeScanner'
            }
            steps{
                echo "Checking code quality"
                withSonarQubeEnv('sonarqube') {
                    sh "${scannerHome}/bin/sonar-scanner"
                }
                timeout(time: 10, unit: 'MINUTES') {
                    waitForQualityGate abortPipeline: true, credentialsId: 'webhook-secret'
                }
            }
            post{
                always{
                    echo "========always========"
                }
                success{
                    echo "========A executed successfully========"
                }
                failure{
                    echo "========A execution failed========"
                }
            }
        }
        stage("Build"){
            steps{
                echo "Building the repository"
                sh 'go build'
                echo "Content in working directory"
                sh "ls -la ."
            }
            post{
                always{
                    echo "========always========"
                }
                success{
                    echo "========A executed successfully========"
                }
                failure{
                    echo "========A execution failed========"
                }
            }
        }
        stage("Run Terraform"){
            steps{
                script {
                    def tfHome = tool name: 'Terraform'
                    env.PATH = "${tfHome}:${env.PATH}"
                }
                
                echo "Running terraform files"
                sh 'terraform --version'
                sh 'terraform init'
            }
            post{
                always{
                    echo "========always========"
                }
                success{
                    echo "========A executed successfully========"
                }
                failure{
                    echo "========A execution failed========"
                }
            }
        }
    }
    post{
        always{
            echo "========always========"
        }
        success{
            echo "========pipeline executed successfully ========"
        }
        failure{
            echo "========pipeline execution failed========"
        }
    }
}