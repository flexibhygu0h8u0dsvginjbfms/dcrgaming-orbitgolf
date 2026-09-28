extends Control
signal scorecard_requested
signal back_requested

# The result of a finished bridge match, from this seat: who won, what the
# stake does, and how far the payout has come. Every word arrives decided from
# the backend's outcome; this screen only lays it out.

const MINT := Color("39f3cd")
const BLUE := Color("79a4ff")
const WHITE := Color("eaf4ff")
const MUTED := Color("92abc2")
const WARN := Color("ffbc62")
const TRACK := Color("26455b")
const INK := Color("071222")

var headline: Label
var score: Label
var money: Label
var note: Label
var steps: Array = []
var step_states: Array = []
var scorecard_button: Button
var back_button: Button

func text(value: String, size_: int, color: Color, parent: Node) -> Label:
	var l := Label.new()
	l.text=value
	l.horizontal_alignment=HORIZONTAL_ALIGNMENT_CENTER
	l.add_theme_font_size_override("font_size",size_)
	l.add_theme_color_override("font_color",color)
	parent.add_child(l)
	return l

func action(value: String, parent: Node, primary := false) -> Button:
	var b := Button.new()
	b.text=value
	b.custom_minimum_size=Vector2(220,48)
	b.add_theme_font_size_override("font_size",17)
	var style := StyleBoxFlat.new()
	style.bg_color=MINT if primary else Color("172e46")
	style.set_corner_radius_all(6)
	b.add_theme_stylebox_override("normal",style)
	var hover: StyleBoxFlat=style.duplicate()
	hover.bg_color=style.bg_color.lightened(0.16)
	b.add_theme_stylebox_override("hover",hover)
	b.add_theme_color_override("font_color",INK if primary else WHITE)
	parent.add_child(b)
	return b

func _ready() -> void:
	set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	mouse_filter=Control.MOUSE_FILTER_STOP
	var shade := ColorRect.new()
	shade.color=Color(0.01,0.025,0.05,0.86)
	shade.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	shade.mouse_filter=Control.MOUSE_FILTER_STOP
	add_child(shade)
	var sheet := PanelContainer.new()
	add_child(sheet)
	sheet.set_anchors_and_offsets_preset(Control.PRESET_CENTER)
	sheet.offset_left=-420
	sheet.offset_right=420
	sheet.offset_top=-270
	sheet.offset_bottom=270
	var style := StyleBoxFlat.new()
	style.bg_color=Color("091727")
	style.border_color=Color("34566f")
	style.set_border_width_all(1)
	style.set_corner_radius_all(14)
	style.content_margin_left=40
	style.content_margin_right=40
	style.content_margin_top=34
	style.content_margin_bottom=34
	sheet.add_theme_stylebox_override("panel",style)
	var column := VBoxContainer.new()
	column.add_theme_constant_override("separation",16)
	sheet.add_child(column)
	text("MATCH OVER",14,MUTED,column)
	headline=text("",64,WHITE,column)
	score=text("",30,WHITE,column)
	money=text("",18,WHITE,column)
	money.autowrap_mode=TextServer.AUTOWRAP_WORD_SMART
	var rail := HBoxContainer.new()
	rail.add_theme_constant_override("separation",12)
	column.add_child(rail)
	for i in range(3):
		var step := VBoxContainer.new()
		step.size_flags_horizontal=Control.SIZE_EXPAND_FILL
		rail.add_child(step)
		var bar := ColorRect.new()
		bar.custom_minimum_size=Vector2(0,4)
		bar.color=TRACK
		step.add_child(bar)
		var name_ := text("",13,MUTED,step)
		steps.append({"bar":bar,"name":name_})
		step_states.append("todo")
	note=text("",15,MUTED,column)
	note.autowrap_mode=TextServer.AUTOWRAP_WORD_SMART
	var spacer := Control.new()
	spacer.size_flags_vertical=Control.SIZE_EXPAND_FILL
	column.add_child(spacer)
	var row := HBoxContainer.new()
	row.alignment=BoxContainer.ALIGNMENT_CENTER
	row.add_theme_constant_override("separation",16)
	column.add_child(row)
	scorecard_button=action("Scorecard",row)
	scorecard_button.pressed.connect(func(): scorecard_requested.emit())
	back_button=action("Back to lobby",row,true)
	back_button.pressed.connect(func(): back_requested.emit())
	hide()

# update_result fills the screen from the backend's outcome.
func update_result(outcome: Dictionary) -> void:
	var kind: String=outcome.get("outcome","")
	headline.text=outcome.get("headline","")
	headline.add_theme_color_override("font_color",MINT if kind=="won" else (BLUE if kind=="lost" else WHITE))
	score.text=outcome.get("score","")
	money.text=outcome.get("money","")
	var payout: Array=outcome.get("payout",[])
	for i in range(steps.size()):
		var step: Dictionary=payout[i] if i<payout.size() else {"label":"","state":"todo"}
		var state: String=step.get("state","todo")
		step_states[i]=state
		steps[i].bar.color=MINT if state=="done" else (WARN if state=="current" else TRACK)
		steps[i].name.text=step.get("label","")
		steps[i].name.add_theme_color_override("font_color",MINT if state=="done" else (WARN if state=="current" else MUTED))
	note.text=outcome.get("payoutNote","")
