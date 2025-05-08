
@echo off
echo Building application...
go build -o .\tmp\main.exe .\cmd\api\main.go
echo Build completed.

@REM Buat file `run.bat`:

@REM ```batch
@REM @echo off
@REM echo Starting application...
@REM hotel-point-app.exe
@REM ```

@REM Buat file `dev.bat`:

@REM ```batch
@REM @echo off
@REM echo Starting in development mode...
@REM go run .\cmd\api\main.go
@REM ```

@REM Buat file `seed.bat`:

@REM ```batch
@REM @echo off
@REM echo Seeding database...
@REM go run .\cmd\seed\main.go
@REM echo Seeding completed.
@REM ```