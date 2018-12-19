//SPDX-License-Identifier: Apache-2.0

var resume = require('./controller.js');

module.exports = function(app){

  app.get('/get_resume/:id', function(req, res){
    resume.get_resume(req, res);
  });
  app.get('/add_resume/:resume', function(req, res){
    resume.add_resume(req, res);
  });
  app.get('/get_all_resumes', function(req, res){
    resume.get_all_resumes(req, res);
  });
}
