'use strict';

angular.module('js', ['ngAnimate', 'ngCookies', 'ngTouch', 'ngSanitize', 'restangular', 'ui.router', 'generator'])
  .config(function($stateProvider, $urlRouterProvider) {
  $stateProvider
    .state('home', {
    url: '/',
    templateUrl: 'app/main/main.html',
    controller: 'MainCtrl'
  })
    .state('generate', {
    url: '/generate',
    templateUrl: 'app/generate/generate.html',
    controller: 'GenerateCtrl'
  })
    .state('send', {
    url: '/send',
    templateUrl: 'app/send/send.html',
    controller: 'SendCtrl'
  });

  $urlRouterProvider.otherwise('/');
});
