'use strict';

angular.module('js')
  .controller('GenerateCtrl', function($scope) {
  $scope.generateForm = {
    length: 12,
    letters: true,
    mixedCase: true,
    numbers: true,
    punctuation: false,
    similarChars: false
  };

});
