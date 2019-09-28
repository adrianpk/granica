var gulp = require("gulp");
var shell = require('gulp-shell');
var exec = require('child_process').exec;

// Compile
gulp.task('install', function (cb) {
  exec('go  build -o ./bin/granica ./cmd/granica.go' , function (err, stdout, stderr) {
    console.log(stdout);
    console.log(stderr);
    cb(err);
  });
})

// Install
gulp.task('restart', function (cb) {
  exec('supervisorctl restart granica', function (err, stdout, stderr) {
    console.log(stdout);
    console.log(stderr);
    cb(err);
  });
})

gulp.task('watch', function() {
  return gulp.watch(["*.go", "./cmd/**/*.go", "./internal/**/*.go", "pkg/**/*.go"], gulp.series('install', 'restart'));
});

gulp.task('default',  gulp.series('watch'));
