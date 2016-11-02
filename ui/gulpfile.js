const gulp = require('gulp');
const browserify = require('browserify');
const babelify = require('babelify');
const source = require('vinyl-source-stream');

const es = require('event-stream');
const concat = require('gulp-concat');
const uglify = require('gulp-uglify');
const minifyCss = require('gulp-minify-css');
const minifyhtml = require('gulp-minify-html');
const sourcemaps = require('gulp-sourcemaps');
const smoosher = require('gulp-smoosher');
const clean = require('gulp-clean');
const rename = require('gulp-rename');
const util = require('util');
const preprocess = require('gulp-preprocess');
 

const qunitVersion = '1.20.0';

gulp.task('css', function() {
    return gulp.src('./src/css/*.css')
        .pipe(minifyCss())
        .pipe(gulp.dest('./dist/css'));
});

gulp.task('browserify', function() {
    var extensions = ['.js', '.json', '.es6', '.jsx'];
    var transforms = ['babelify', {
        presets: ["es2015"]
    }];

    var browserifyTask = browserify({
            entries: './src/js/index.jsx',
            extensions: extensions,
            debug: true,
        })
        .transform(transforms)
        .bundle()
        .pipe(source('main.js'))
        .pipe(gulp.dest('dist/js'));

    var copyTestFiles = gulp.src(util.format('src/js/qunit-%s.js', qunitVersion)).pipe(gulp.dest('dist/js'));

    return es.merge(browserifyTask, copyTestFiles);
});

gulp.task('inline', ['browserify', 'css'], function() {
    var main = gulp.src('src/index.html')
        .pipe(preprocess({context: { DEBUG: false}}))
        .pipe(smoosher({
            base: 'dist'
        }))
        .pipe(minifyhtml())
        .pipe(gulp.dest('dist'));

    var mainTest = gulp.src('src/index.html')
        .pipe(preprocess({context: { DEBUG: true}}))
        .pipe(smoosher({
            base: 'dist'
        }))
        .pipe(minifyhtml())
        .pipe(rename('index-test.html'))
        .pipe(gulp.dest('dist'));

    return es.merge(main, mainTest);
});

gulp.task('copyToPublic', ['inline'], function() {
    return es.merge(
        gulp.src(['./dist/index.html']).pipe(gulp.dest('./public')),
        gulp.src(['./dist/index-test.html']).pipe(gulp.dest('./public')),
        gulp.src(['src/fonts/**/*']).pipe(gulp.dest('./public/fonts'))
    )
});

gulp.task('cleanup', ['inline', 'copyToPublic'], function() {
    gulp.src('dist/js', {
        read: false
    }).pipe(clean())
    return gulp.src('dist/css', {
        read: false
    }).pipe(clean())
});

gulp.task('default', ['css', 'browserify', 'inline', 'cleanup']);

gulp.task('watch', function() {
    gulp.watch('src/**/*.*', ['default']);
});
