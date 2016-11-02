package cmd

import (
	"fmt"
	"math"

	"github.com/fatih/color"
)

//LogoProgress ...
type LogoProgress struct {
	lines []string
	last  int
}

//NewLogoProgress ...
func NewLogoProgress() *LogoProgress {
	var lines []string

	red := color.New(color.FgRed, color.Bold).SprintFunc()
	white := color.New(color.FgWhite, color.Bold).SprintFunc()

	lines = append(lines,
		fmt.Sprintf("                      %s                     \n", white(",▄▄▄▄,")),
		fmt.Sprintf("                   %s  %s%s                 \n", white("▄▓█▀Γ"), red("."), white("╙▀▓▄,")),
		fmt.Sprintf("               %s  %s   %s             \n", white("▄▓▓▀Γ"), red("-T░U░="), white("▀▀▓▄,")),
		fmt.Sprintf("           %s  %s%s%s  %s          \n", white("╓▄▓▀▀"), red("-=░UUU"), white("╓▄"), red("UUUU░¬"), white(".▀█▓▄")),
		fmt.Sprintf("       %s  %s%s%s  %s      \n", white(",▄▓▀▀"), red(".=░UUU"), white("╓▄▄▓▓▓▓▓▄▄"), red("UUUUS-"), white("└▀█▓▄")),
		fmt.Sprintf("     %s %s%s%s  %s   \n", white(",▓▀."), red(".¬░UUUU"), white("▄▄▓▓▓▓▓▓▓▓▓▓▓▓▓▄▄"), red("UUU░S-"), white("'▀▓▄")),
		fmt.Sprintf("    %s  %s%s%s%s%s  %s  \n", white("▓▀"), red("SUUUU"), white("▄▄▓▓▓▓▓▓▓▀▀"), red("÷⌂Γ"), white("▀█▓▓▓▓▓▓"), red("UUUUUU░="), white("▓▌")),
		fmt.Sprintf("   %s  %s%s%s%s%s %s  \n", white("▐▓"), red("UUU"), white("j▓▓▓▓▓▓▓█▀"), red("⌂+UUUUUUU-Γ"), white("▀▓▓▓"), red("UUUUUUUUL"), white("▐▓")),
		fmt.Sprintf("   %s  %s%s%s%s%s %s  \n", white("▐▓"), red("UUU"), white("▐▓▓▓▓▀´"), red("⌂UUUUUUUUUUUUUUU"), white("Ÿ╙"), red("UUUUUUUUL"), white("▐▓")),
		fmt.Sprintf("   %s  %s%s%s %s   \n", white("▐▓"), red("UUU"), white("▐▓▓▓▓"), red("UUUUUUUUUUUUUUUUUUUUUUUUUUUUL"), white("▐▓")),
		fmt.Sprintf("   %s  %s%s%s %s   \n", white("▐▓"), red("UUU"), white("▐▓▓▓▓"), red("UUUUUUUUUUUUUUUUUUUUUUUUUUUUL"), white("▐▓")),
		fmt.Sprintf("   %s  %s%s%s %s   \n", white("▐▓"), red("UUU"), white("▐▓▓▓▓"), red("UUUUUUUUUUUUUUUUUUUUUUUUUUUUL"), white("▐▓")),
		fmt.Sprintf("   %s  %s%s%s %s   \n", white("▐▓"), red("UUU"), white("▐▓▓▓▓"), red("UUUUUUUUUUUUUUUUUUUUUUUUUUUUL"), white("▐▓")),
		fmt.Sprintf("   %s  %s%s%s %s   \n", white("▐▓"), red("UUU"), white("▐▓▓▓▓"), red("UUUUUUUUUUUUUUUUUUUUUUUUUUUUL"), white("▐▓")),
		fmt.Sprintf("   %s  %s%s%s%s%s %s   \n", white("▐▓"), red("UUU"), white("▐▓▓▓▓▓▄▄"), red("+UUUUUUUUUUUUU"), white(",▄▓"), red("UUUUUUUUL"), white("▐▓")),
		fmt.Sprintf("   %s  %s%s%s%s%s %s   \n", white("▐▓"), red("UUU⌂"), white("╙▀▓▓▓▓▓▓▓▄▄"), red("UUUUUU"), white(",▄▄▓▓▓▓"), red("UUUUUUUU¬"), white("▐▓")),
		fmt.Sprintf("    %s  %s%s%s %s   \n", white("▓▓"), red("¬░UUUU"), white("`▀▀▓▓▓▓▓▓▓▄▄▄▓▓▓▓▓▓▓▀"), red("UUUUUU░'"), white(",▓▀")),
		fmt.Sprintf("     %s  %s%s%s  %s    \n", white("▀▓▓▄"), red("'T░UUU+"), white("`▀█▓▓▓▓▓▓▓▓▓▀▀"), red("⌂UUUU░'"), white(",▄▓█´")),
		fmt.Sprintf("        %s  %s%s%s   %s       \n", white("▀▀▓▄,"), red("'TUUUU⌂"), white("``▀█▓▀▀"), red("÷UUUU░¬"), white("▄▓█▀´")),
		fmt.Sprintf("            %s  %s   %s           \n", white("▀▀▓▄,"), red("'░UUUUUUUU░¬"), white("▄▄▓▀▀")),
		fmt.Sprintf("               %s   %s  %s               \n", white("`▀█▓▄"), red("'T░E'"), white("╓▄▓▀▀")),
		fmt.Sprintf("                   %s                  \n", white("└▀█▓▄▄╓▄▄▓▀▀´")),
		fmt.Sprintf("                       %s                       \n", white("`▀▀´")),
		fmt.Sprint("\n"),
	)

	return &LogoProgress{lines, -1}
}

//Set ...
func (b *LogoProgress) Set(progress int) error {
	a := float64(progress) / 100 * float64(len(b.lines))
	index := int(math.Floor(a))
	//fmt.Printf("index: %v last: %v\n", index, b.last)
	if index >= len(b.lines) {
		return nil
	}
	if index > b.last {
		for b.last < index {
			b.last++
			fmt.Fprint(color.Output, b.lines[b.last])
		}
	}
	return nil
}
