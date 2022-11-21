
![Purple__Yellow_Career_Coach_LinkedIn_Banner(1)](https://user-images.githubusercontent.com/60812924/201931547-8aa0bb70-5510-4f27-9582-d044162089e5.png)

### 💫About Kritagya

---

Layoffs are a regular corporate response to the difficulties created by technological advancements and escalating global competitiveness.Developers and businesses can interact and share information via Slack.

Recently, it has become more and more common for engineers to manage their small team of engineers.

In order to manage Jenkins build results, incorporate Pull requests, and notify the team of open issues using [Github Actions](https://github.com/archive/github-actions-slack), Kritagya may link to the Slack workspace and channels. To send requests to a Go backend, the bot programme will use Websocket, also referred to as "Socket-Mode" in the Slack community.

💫 **How to set up Kritagya locally** 

---

1. Set up a [new slack app](https://api.slack.com/authentication/basics#creating) for your slack channel using  the necessary scopes and slash commands [used](https://gist.github.com/ShivamTyagi12345/419d2319674fa8cabb369482470565e3) in Kritagya.
2. On a local computer where you have a Go installed, clone this repository and go to `.env.example`  file. Copy the contents
3. Next, create a new file titled `.env` and paste the contents ,followed by replacing the values available in [slack-api](https://api.slack.com/apps) dashboard 
4. Go to *Jenkinsfile* and replace *credentialsId*  with your ssh key and *url* with your repository url (*see*,line 8) , Next, replace the *slack-channel* in  repository.yml with the channelID ,you want to send notifications to                        
5. Run `go run main.go https://localhost:3000 SUCCESS test 1` in terminal to run a dummy test 🎊

💫**How to use commands of Kritagya**

---

1.  ***hello @Kritagya*** : This command Greets you with your username and ensures that you have a lovely day

![Untitled](https://user-images.githubusercontent.com/60812924/201941587-6ef4e782-ba96-4dd2-8bab-1af481b36f26.png)

2.  ***<Any random text followed by> @Kritagya :*** this helps us know that you are ready to throw instructions at us ❤️
******

![Untitled](Public/Untitled%201.png)

3. ***/namaste :***  ******[******Optional slash commands******]****** now that we have our bot ready , lets learn some French using */namaste* 

![Untitled](Public/Untitled%202.png)

4. Now create/close any issue in your repository, and notice how Kritagya’s automated tools pick it up **instantly** to flashes us the message 

![image](https://user-images.githubusercontent.com/60812924/201945629-db9c6acc-7432-4f19-acb7-b4f49baaefe6.png)


- *Note*: 

        👉 This tells us which branch has opened or closed  an issue  (here, main )
        👉 This will also tell us about the URL  

5. *Now make some changes in your local copy of Kritagya and push the changes to `main` branch* 

![WhatsApp Image 2022-11-21 at 20 41 09](https://user-images.githubusercontent.com/94890149/203089843-397b8a11-91f0-49ed-ad03-d0b41260f17b.jpeg)


- *Note*: 

        👉 This tells us which branch has new commits pushed  (here, main )
        👉 This will also tell us about the URL  


6. *Start Your Jenkins Pipeline and run a build*

Notice how almost Instantly a pipeline is finished and Krityaga Notifies it with either a not success or Success result.

![Untitled](Public/Untitled%205.png)

Successful with URl of the Jenkins build

7. Merge a Pull Request onto ***main*** branch :  Welcome new Pull Requests, Kritagya makes sure that no PR goes skips the maintainers eyes 👏

![Untitled](Public/Untitled%208.png)


💫 **Scalability Scope**

---

The slackbot programme must be updated when new versions of the slackbot programme are available if you decide to package the slackbot software into an executable that can be deployed on the Jenkins instance. When the Jenkins pipeline script runs, you must make sure that the go programming language is installed and that you can clone down the most recent version of the go application from github.

💫**Privacy and Security**

---

Currently , one vulnerability with Kritagya happens to be using the https://github.com/archive/github-actions-slack 
This essestialy doesn’t accept environment variables as the slack channel-ID, This can be removed by upgrading to the paid version available, or building a environment configuration management tool from scratch

💫 **Challenges**

---

We can change our turning a slackbot software into microservice utilise just one API endpoint endpoint requires a user-hash in json. the payload in json and the payload What were the previous orders? line references if we want to run the Slackbot application. Is this the only dependence for Jenkins? needs Does the curl tool allow you to call UH? the API's json payload endpoint There will be no need to leave as a result. programming language installed on jenkins or by setting up the Slackbot programme itself on Jenkins Additionally, this will generally be a scalable substitute for what we had earlier
