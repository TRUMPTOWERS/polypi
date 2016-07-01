var gulp = require('gulp');

gulp.task('copy', function() {
    return gulp.src([
        './bower_components/webcomponentsjs/webcomponents.js',
        './bower_components/font-awesome/',
        './src/logo.svg'
    ])
    .pipe(gulp.dest('./dist/static'));
});

gulp.task('default', ['copy']);
