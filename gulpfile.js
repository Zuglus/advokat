// 'use strict'
const gulp = require('gulp');
const scss = require('gulp-sass');
const autoprefixer = require('gulp-autoprefixer');
const browserSync = require('browser-sync');
const pathSCSS = 'app/scss/**/*.scss';
const pathHTML = 'app/*.html';



gulp.task('scss', function() {
  return gulp.src(pathSCSS)
  .pipe(scss({ outputStyle: 'compressed' }))
  .pipe(autoprefixer({
    cascade: false
  }))
  .pipe(gulp.dest('app/css'))
  .pipe(browserSync.reload({stream: true}));
});

gulp.task('html', function() {
  return gulp.src(pathHTML)
  .pipe(browserSync.reload({stream: true}))
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
  gulp.watch(pathHTML, gulp.parallel('html'))
});

gulp.task('default', gulp.parallel('browser-sync', 'watch'));