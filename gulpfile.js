// 'use strict'
const gulp = require('gulp'),
scss = require('gulp-sass'),
browserSync = require('browser-sync');

const pathSCSS = 'app/scss/**/*.scss'

gulp.task('scss', function() {
  return gulp.src(pathSCSS)
  .pipe(scss({ outputStyle: 'expanded' }))
  .pipe(gulp.dest('app/css'))
  .pipe(browserSync.reload({stream: true}));
});

gulp.task('browser-sync', function() {
  browserSync.init({
    server: {
      baseDir: "app/"
    }
  });
});

gulp.task('watch', function() {
  gulp.watch(pathSCSS, gulp.parallel('scss'))
});

gulp.task('default', gulp.parallel('browser-sync', 'watch'));