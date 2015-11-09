angular.module("Namely_Dev_Portal", ['mgcrea.ngStrap', 'ngAnimate', 'ngTable', 'ui.router', 'LocalStorageModule', 'angular-jwt'])
      .config(function(localStorageServiceProvider){
        localStorageServiceProvider
          .setPrefix('Namely_Developers')
          .setStorageType('sessionStorage');
      })
      .config(function($stateProvider, $urlRouterProvider) {
        $urlRouterProvider.otherwise("/");
        $stateProvider
          .state('login', {
            url: "/",
            templateUrl: "/js/_login.html",
            controller: "LoginCtrl"
          })
          .state('begin_reset_password', {
            url: "/password",
            templateUrl: "/js/_password_reset.html",
            controller: "PasswordResetCtrl"
          })
          .state('complete_reset_password', {
            url: "/password/:token",
            templateUrl: "/js/_password_reset.html",
            controller: "PasswordResetCtrl"
          })
          .state('main', {
            abstract: true,
            url: "/main",
            templateUrl: "/js/_main.html",
            controller: "MainCtrl"
          })
          .state('main.admin', {
            abstract: true,
            url: "/admin",
            templateUrl: "/js/admin/_admin_main.html",
            controller: "AdminCtrl"
          })
          .state('main.admin.users', {
            url: "/users",
            templateUrl: "/js/admin/_users.html",
            controller: "UserCtrl"
          })
          .state('main.admin.people', {
            url: "/people",
            templateUrl: "/js/admin/_people.html",
            controller: "PeopleCtrl"
          })
          .state('main.dev', {
            abstract: true,
            url: "/developer",
            templateUrl: "/js/developers/_developers_main.html",
            controller: "DevCtrl"
          })
          .state('main.dev.profile', {
            url: "/profile",
            templateUrl: "/js/developers/_profile.html",
            controller: "DevProfileCtrl"
          })
          .state('main.dev.apps', {
            url: "/apps",
            templateUrl: "/js/developers/_apps.html",
            controller: "DevAppCtrl"
          })
      })
      .config(function($httpProvider, jwtInterceptorProvider) {
        jwtInterceptorProvider.tokenGetter = ['localStorageService', function(localStorageService) {
          return localStorageService.get('token');
        }];
        $httpProvider.interceptors.push('jwtInterceptor');
        $httpProvider.interceptors.push(['$q', '$injector', function($q, $injector){
          var AuthenticationInterceptor = {
            responseError: function(response){
              AuthService = $injector.get("AuthService")
              $http = $injector.get("$http")
              if(!AuthService.refreshToken){
                AuthService.logout();
                return $q.reject(response)
              }
              if(response.status == 401){
                var deferred = $q.defer()
                AuthService.refresh().then(
                function(data){
                  AuthService.registerToken(data.data.token)
                  deferred.resolve()
                },
                function(data){
                  AuthService.logout();
                  deferred.reject()
                })
                return deferred.promise.then(function(){
                  response.config.headers['Authorization'] = "Bearer " + AuthService.token
                  return $http(response.config)
                })
              }
              return $q.reject(response)
            }
          }
          return AuthenticationInterceptor
        }])
      })
      .service('AuthService', ['$http','$q', '$state', 'localStorageService', '$alert', 'jwtHelper', function($http, $q, $state, localStorageService, $alert, jwtHelper){
        this.loggedIn = false;
        this.token = null;
        this.user = null;
        this.exp = null;
        this.refreshToken= null;

        this.routeLogin=function(admin){
          if (admin) {
            $state.go('main.admin.users');
          }else if (!admin) {
            $state.go('main.dev.apps');
          }else {
            $state.go('login');
          }
        }

        this.GetUsername = function(){
          if(this.user){
            return user.email;
          }
        }

        this.refresh = function(){
          if(this.refreshToken){
            return $http.post('/refresh', {"token": this.refreshToken})
          }else{
            $state.go('login');
          }
        }

        this.logout=function(){
          $http.get('/logout')
          localStorageService.clearAll();
          this.token = null;
          this.user = null;
          this.loggedIn = false;
          $state.go('login')
        }

        this.isLoggedIn=function(){
          if(this.user){
            if(this.exp < Date.now()){
              this.loggedIn = true;
              return true;
            }else{
              return false;
            }
          }else{
            this.user = localStorageService.get('user');
            if(this.user){
              this.exp = localStorageService.get('token_exp');
              this.refreshToken= localStorageService.get('refresh_token');
              if(this.exp < Date.now()){
                this.loggedIn = true;
                return true;
              }else{
                this.loggedIn = false;
                this.logout();
              }
            }else{
              this.loggedIn = false;
              return false;
            }
          }
        }

        this.registerToken=function(token){
          this.token = token;
          tokenpayload = jwtHelper.decodeToken(this.token);
          this.exp = tokenpayload.exp
          this.user = angular.fromJson(tokenpayload.usr)
          localStorageService.set('user', this.user)
          localStorageService.set('token', this.token)
          localStorageService.set('token_exp', this.exp)
        }

        this.login=function(username, password){
          var deferred = $q.defer();
          $http
            .post('/login', {"username": username, "password": password})
            .success(function(data){
              if(!data.hasOwnProperty('status')){
                this.loggedIn = true;
                this.refresh_token= data.refresh;
                localStorageService.set('refresh_token', this.refresh_token)
                this.token = data.token;
                tokenpayload = jwtHelper.decodeToken(this.token);
                this.exp = tokenpayload.exp
                this.user = angular.fromJson(tokenpayload.usr)
                localStorageService.set('user', this.user)
                localStorageService.set('token', this.token)
                localStorageService.set('refresh_token', this.refresh_token)
                localStorageService.set('token_exp', this.exp)
                deferred.resolve(this.user);
              }else{
                deferred.resolve(data);
              }
            })
            .error(function(data, status){
              if(status==404){
                $alert({title:"Login Failure", placement: "bottom", content:"Invalid Username/Password", show: true, type:"warning", duration: "3", container: '#fullscreen_bg', animation: 'am-slide-top'});
              }else{
                $alert({title:"Server Error", placement: "bottom", content:"Log in unavailable.  Try again later" , show: true, type:"danger", duration: "3", container: '#fullscreen_bg', animation: 'am-slide-top'});
              }
            });
          return deferred.promise;
        }
      }])
      .factory('UsersFactory', function ContentFactory($http){
        return {
          getUsers: function(){
            return $http.get('/v1/users');
          },
          addUser: function(newUser){
            return $http.post('/v1/users', newUser);
          },
          updateUser: function(update){
            return $http.put('/v1/users', update);
          },
          deleteUser: function(user_id){
            return $http.delete('/v1/users/'+ user_id);
          },
          sendResetRequest: function(email){
            return $http.get('/reset_password/'+email);
          },
          completeReset: function(request){
            return $http.post('/complete_reset/', request);
          }
        }
      })
      .factory('PeopleFactory', function ContentFactory($http){
        return {
          getPeople: function(){
            return $http.get('/v1/people');
          },
          getPerson: function(people_id){
            return $http.get('/v1/person/' + app_id);
          },
          /*
          getUserApps: function(people_id){
            return $http.get('/v1/users/' + user_id + '/apps');
          },
          */
          addPerson: function(newPeople){
            return $http.post('/v1/people', newPeople);
          },
          deleteApp: function(app_id){
            return $http.delete('/v1/apps/'+ app_id);
          }
        }
      })
      .controller('LoginCtrl', ['$scope', '$state', 'AuthService', '$alert', function($scope, $state, AuthService, $alert){
        $scope.credentials={};
        var init = function(){
          if(AuthService.isLoggedIn()){
            AuthService.routeLogin(AuthService.user.admin); 
          }
        }

        $scope.processLogin = function(){
          AuthService.login($scope.credentials.username, $scope.credentials.password)
          .then(function(payload){
            AuthService.routeLogin(payload.admin); 
          });
          $scope.credentials.$pristine=true;
          $scope.credentials.username="";
          $scope.credentials.password="";
        }
        init();
      }])
      .controller('MainCtrl', ['$scope', 'AuthService', function($scope, AuthService){
        if(!AuthService.isLoggedIn()){
          AuthService.logout();
          return
        }
        $scope.user_name = AuthService.user.first_name + " " + AuthService.user.last_name;
      }])
      .controller('AdminCtrl', ['$scope', 'AuthService', '$state', function($scope, AuthService, $state){
        $scope.user_name = AuthService.user.first_name + " " + AuthService.user.last_name;
        var init = function(){
          if(!AuthService.isLoggedIn()){
            $state.go('login');
          }
        }
        $scope.logout=function(){
          AuthService.logout();
        }
        init();
      }])
      .controller('PasswordResetCtrl', ['$scope', 'AuthService', '$stateParams', 'UsersFactory', '$alert', function($scope, AuthService, $stateParams, UsersFactory, $alert){
        $scope.entry_state = true;
        var request= {};
        $scope.passwords = {};
        $scope.reset = {};
        $scope.passwordMatch = false;
        if ($stateParams.token == "" || !$stateParams.hasOwnProperty("token")) {
          $scope.entry_state = true;
        }else{
          $scope.entry_state = false;
        }

        $scope.completeReset = function(){
          //Verify that password fields match
          if($scope.passwords.first_time == $scope.passwords.second_time){
            request.password= $scope.passwords.first_time
            request.token = $stateParams.token
            UsersFactory.completeReset(request)
            .success(function(data){
              $alert({title:"Password Changed", placement: "top", content:"Success!", show: true, type:"success", container: 'body', animation: 'am-slide-top', duration: '4'});
              //Expose login link
              AuthService.logout();
            })
            .error(function(data, status){
              if(status == "404"){
                $alert({title:"Whups", placement: "top", content:"Provided token is invalid.  Double check the link in your email.", show: true, type:"warning", container: 'body', animation: 'am-slide-top', duration: '4'});
              }else{
                $alert({title:"Whups", placement: "top", content:"Something went wrong.  Please contact administrator.", show: true, type:"warning", container: 'body', animation: 'am-slide-top', duration: '4'});
              }
            })
            return
          }else{
            $alert({title:"Error ", placement: "top", content:"Passwords must match", show: true, type:"warning", container: 'body', animation: 'am-slide-top', duration: '4'});
          }
        }

        $scope.beginReset = function(){
          UsersFactory.sendResetRequest($scope.reset.email)
          .success(function(data, status, codes, config){
            $alert({title:"Email Sent. ", placement: "top", content:"Go check your email.", show: true, type:"success", container: 'body', animation: 'am-slide-top', duration: '4'});
            $scope.reset.email = "";
          })
          .error(function(data, status, codes, config){
            if(status == "404"){
              $alert({title:"Error ", placement: "top", content:"Email address not found.", show: true, type:"warning", container: 'body', animation: 'am-slide-top', duration: '4'});
            }else{
              $alert({title:"Error ", placement: "top", content:"Please contact administrator", show: true, type:"warning", container: 'body', animation: 'am-slide-top', duration: '4'});
            }
          })
        }
      }])
      .controller('AdminCtrl', ['$scope', 'AuthService', '$state', function($scope, AuthService, $state){
        var init = function(){
          if(!AuthService.isLoggedIn()){
            $state.go('login');
          }
        }
        $scope.logout=function(){
          AuthService.logout();
        }
        init();
      }])
      .controller('DevCtrl', ['$scope', 'AuthService', '$state', function($scope, AuthService, $state){
        var init = function(){
          if(!AuthService.isLoggedIn()){
            $state.go('login');
          }
        }
        $scope.logout=function(){
          AuthService.logout();
        }
        init();
      }])
      .controller('UserCtrl', ['$scope', 'ngTableParams', 'UsersFactory', '$window', '$alert', 'AuthService', function($scope, ngTableParams, UsersFactory, $window, $alert, AuthService){
        var init = function(){
          if(!AuthService.isLoggedIn()){
            $state.go('login');
          }
        }
        init();
        $scope.editId = -1;
        $scope.users=[];
        $scope.clearForms=function(){
          $scope.new_user={};
          $('#new_user').modal('hide');
        }
        $scope.processUser = function(){
          UsersFactory.addUser($scope.new_user).success(function(data,status){
            $scope.users.push(data);
            $scope.usertableParams.total($scope.users.length);
            $scope.usertableParams.reload();
          });
          $scope.clearForms();
        }
        UsersFactory.getUsers().success(function(data, status){
          if(data.status!="No Users found"){
            $scope.users=data
            $scope.usertableParams.total($scope.users.length);
            $scope.usertableParams.reload();
          }else{
            $scope.users=[];
          }
        })
        $scope.usertableParams = new ngTableParams({
            page:1,
            count: 10
          },
          {
            total: 0,
            getData: function($defer, params){
              var data = $scope.users;
              if(data){
                $defer.resolve(data.slice((params.page() - 1) * params.count(), params.page() * params.count()));
              }
            }
        });

        $scope.setEditId = function(id){
          $scope.editId = id;
        }

        $scope.closeUserChanges = function(){
          $scope.editId = -1;
        }

        $scope.deleteUser = function(entity){
          if($window.confirm("Are you sure you wish to delete this user?")){
            UsersFactory.deleteUser(entity.id).success(function(data, status){
              $alert({title:"Success", placement: "top", content:"User Deleted", show: true, type:"success", duration: "3", container: 'body', animation: 'am-slide-top'});
              $scope.closeUserChanges();
              i = $scope.users.indexOf(entity);
              $scope.users.splice(i, 1)
              $scope.usertableParams.total($scope.users.length);
              $scope.usertableParams.reload();
            });
          }
        }

        $scope.saveUserChanges = function(entity){
          //TODO: Don't update table entity on failure
          UsersFactory.updateUser(entity).success(function(data, status){
            $alert({title:"Success", placement: "top", content:"User updated successfully", show: true, type:"success", duration: "3", container: 'body', animation: 'am-slide-top'});
            $scope.closeUserChanges()
          });
        }
      }])
      .controller('PeopleCtrl', ['$scope', 'ngTableParams', 'PeopleFactory', '$alert', 'AuthService', function($scope, ngTableParams, PeopleFactory, $alert, AuthService){
        var init = function(){
          if(!AuthService.isLoggedIn()){
            $state.go('login');
          }
        }
        $scope.people= [];
        $scope.editId = -1;
        $scope.clearForms=function(){
          $scope.new_people={};
          $('#new_people').modal('hide');
        }
        PeopleFactory.getPeople().success(function(data, status){
          if(data.status!="No People found"){
            $scope.people=data
            $scope.peopletableParams.total($scope.people.length);
            $scope.peopletableParams.reload();
          }else{
            $scope.people=[];
          }
        })
        $scope.peopletableParams= new ngTableParams({
            page:1,
            count: 10
          },
          {
            total: 0,
            getData: function($defer, params){
              var data = $scope.apps;
              if(data){
                $defer.resolve(data.slice((params.page() - 1) * params.count(), params.page() * params.count()));
              }
            }
        });
        $scope.setEditId = function(id){
          $scope.editId = id;
        }

        $scope.closePeopleChanges= function(){
          $scope.editId = -1;
        }

        $scope.deletePerson = function(entity){
          PeopleFactory.deleteApp(entity.client_id).success(function(data, status){
            $alert({title:"Success", placement: "top", content:"App Deleted", show: true, type:"success", duration: "3", container: 'body', animation: 'am-slide-top'});
            $scope.closeAppChanges();
            i = $scope.apps.indexOf(entity);
            $scope.apps.splice(i, 1)
            $scope.apptableParams.total($scope.apps.length);
            $scope.apptableParams.reload();
          });
        }

        $scope.saveAppChanges = function(entity){
          //TODO: Don't update table entity on failure
          PeopleFactory.updateApp(entity).success(function(data, status){
            $alert({title:"Success", placement: "top", content:"App updated successfully", show: true, type:"success", duration: "3", container: 'body', animation: 'am-slide-top'});
            $scope.closeAppChanges()
          });
        }
        init();
      }])
      .controller('DevAppCtrl', ['$scope', 'ngTableParams', 'PeopleFactory', '$alert', 'AuthService', function($scope, ngTableParams, PeopleFactory, $alert, AuthService){
        $scope.apps= [];
        $scope.editId = -1;
        $scope.clearForms=function(){
          $scope.new_app={};
          $('#new_app').modal('hide');
        }
        PeopleFactory.getUserApps(AuthService.user.id).success(function(data, status){
          if(status==200){
            $scope.apps=data
            $scope.apptableParams.total($scope.apps.length);
            $scope.apptableParams.reload();
          }else{
            $scope.apps=[];
          }
        })
        $scope.apptableParams= new ngTableParams({
            page:1,
            count: 10
          },
          {
            total: 0,
            getData: function($defer, params){
              var data = $scope.apps;
              if(data){
                $defer.resolve(data.slice((params.page() - 1) * params.count(), params.page() * params.count()));
              }
            }
        });
        $scope.setEditId = function(id){
          $scope.editId = id;
        }

        $scope.closeAppChanges= function(){
          $scope.editId = -1;
        }

        $scope.deleteApp = function(entity){
          PeopleFactory.deleteApp(entity.client_id).success(function(data, status){
            $alert({title:"Success", placement: "top", content:"App Deleted", show: true, type:"success", duration: "3", container: 'body', animation: 'am-slide-top'});
            $scope.closeAppChanges();
            i = $scope.apps.indexOf(entity);
            $scope.apps.splice(i, 1)
            $scope.apptableParams.total($scope.apps.length);
            $scope.apptableParams.reload();
          });
        }

        $scope.saveAppChanges = function(entity){
          //TODO: Don't update table entity on failure
          PeopleFactory.updateApp(entity).success(function(data, status){
            $alert({title:"Success", placement: "top", content:"App updated successfully", show: true, type:"success", duration: "3", container: 'body', animation: 'am-slide-top'});
            $scope.closeAppChanges()
          });
        }

        $scope.processApp = function(){
          $scope.new_app.user_id = AuthService.user.id;
          PeopleFactory.addApp($scope.new_app).success(function(data,status){
            $scope.apps.push(data);
            $scope.apptableParams.total($scope.apps.length);
            $scope.apptableParams.reload();
          });
          $scope.clearForms();
        }

        $scope.refreshAppSecret = function(entity, index){
          //Call PeopleFactory APi to get new thing
          PeopleFactory.getNewAppSecret(entity)
          .success(function(response){
            $scope.apps[index].secret = response.secret
          })
          .error(function(response, status){
            $alert({title:"Failure", placement: "top", content:"Unable to reset Secret", show: true, type:"danger", duration: "3", container: 'body', animation: 'am-slide-top'});
          })
        }
      }])
      .controller('DevProfileCtrl', ['$scope', 'AuthService', 'UsersFactory', '$alert', function($scope, AuthService, UsersFactory, $alert){
        $scope.user = AuthService.user;
        $scope.saveProfileChanges = function(user) {
          //Take form fields and save them to the server  
          UsersFactory.updateUser(user)
          .success(function(response){
            $alert({title:"Success", placement: "top", content:"Profile Updated", show: true, type:"success", duration: "4", container: 'body', animation: 'am-slide-top'});
          })
          .error(function(response, status){
            $alert({title:"Failure", placement: "top", content:"Profile could not be updated.", show: true, type:"danger", duration: "4", container: 'body', animation: 'am-slide-top'});
          })
        }
      }])
