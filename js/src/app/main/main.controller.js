'use strict';

angular.module('js')
  .controller('MainCtrl', function ($scope, $location) {
    $scope.generate = function() {
      $location.path('generate');
    };
    $scope.send = function() {
      $location.path('send');
    };
  });
