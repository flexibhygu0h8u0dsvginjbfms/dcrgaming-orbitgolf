extends Control
signal back_requested
signal fund_requested

# The table lobby: a table from its accepted invite to the first shot. Every
# word and state arrives decided from the backend's lobby description; this
# screen only lays it out.

const MINT := Color("39f3cd")
const BLUE := Color("79a4ff")
const WHITE := Color("eaf4ff")
const MUTED := Color("92abc2")
const WARN := Color("ffbc62")
const WRONG := Color("ff7a7a")
const TRACK := Color("26455b")
const INK := Color("071222")

var table_label: Label
var steps: Array = []
var step_states: Array = []
var seat_cards: Array = []
var terms: Dictionary = {}
var next_step: Label
var next_detail: Label
var error_label: Label
var fund_button: Button
var back_button: Button

func box(parent: Node, bg: Color) -> VBoxContainer:
	var p := PanelContainer.new()
	var style := StyleBoxFlat.new()
	style.bg_color=bg
	style.border_color=Color("34566f")
	style.set_border_width_all(1)
	style.set_corner_radius_all(12)
	style.content_margin_left=20
	style.content_margin_right=20
	style.content_margin_top=16
	style.content_margin_bottom=16
	p.add_theme_stylebox_override("panel",style)
	parent.add_child(p)
	var v := VBoxContainer.new()
	v.add_theme_constant_override("separation",8)
	p.add_child(v)
	return v

func text(value: String, size_: int, color: Color, parent: Node) -> Label:
	var l := Label.new()
	l.text=value
	l.add_theme_font_size_override("font_size",size_)
	l.add_theme_color_override("font_color",color)
	parent.add_child(l)
	return l

func action(value: String, parent: Node, primary := false) -> Button:
	var b := Button.new()
	b.text=value
	b.custom_minimum_size=Vector2(0,44)
	b.add_theme_font_size_override("font_size",16)
	var style := StyleBoxFlat.new()
	style.bg_color=MINT if primary else Color("172e46")
	style.set_corner_radius_all(6)
	style.content_margin_left=16
	style.content_margin_right=16
	b.add_theme_stylebox_override("normal",style)
	var hover: StyleBoxFlat=style.duplicate()
	hover.bg_color=style.bg_color.lightened(0.16)
	b.add_theme_stylebox_override("hover",hover)
	var off: StyleBoxFlat=style.duplicate()
	off.bg_color=Color("172e46")
	b.add_theme_stylebox_override("disabled",off)
	b.add_theme_color_override("font_color",INK if primary else WHITE)
	b.add_theme_color_override("font_disabled_color",MUTED)
	parent.add_child(b)
	return b

func _ready() -> void:
	set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	mouse_filter=Control.MOUSE_FILTER_STOP
	var shade := ColorRect.new()
	shade.color=Color(0.01,0.025,0.05,0.92)
	shade.set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	shade.mouse_filter=Control.MOUSE_FILTER_STOP
	add_child(shade)
	var root := VBoxContainer.new()
	root.position=Vector2(48,36)
	root.custom_minimum_size=Vector2(1344,828)
	root.add_theme_constant_override("separation",18)
	add_child(root)

	var header := HBoxContainer.new()
	header.add_theme_constant_override("separation",16)
	root.add_child(header)
	back_button=action("<  BACK TO LOBBY",header)
	back_button.pressed.connect(func(): back_requested.emit())
	var titles := VBoxContainer.new()
	titles.size_flags_horizontal=Control.SIZE_EXPAND_FILL
	header.add_child(titles)
	text("TABLE LOBBY",30,WHITE,titles)
	table_label=text("",13,MUTED,titles)

	var rail := HBoxContainer.new()
	rail.add_theme_constant_override("separation",10)
	root.add_child(rail)
	for i in range(6):
		var step := VBoxContainer.new()
		step.size_flags_horizontal=Control.SIZE_EXPAND_FILL
		rail.add_child(step)
		var bar := ColorRect.new()
		bar.custom_minimum_size=Vector2(0,4)
		bar.color=TRACK
		step.add_child(bar)
		var name_ := text("%d" % (i+1),13,MUTED,step)
		steps.append({"bar":bar,"name":name_})
		step_states.append("todo")

	var middle := HBoxContainer.new()
	middle.add_theme_constant_override("separation",18)
	middle.size_flags_vertical=Control.SIZE_EXPAND_FILL
	root.add_child(middle)
	var left := VBoxContainer.new()
	left.size_flags_horizontal=Control.SIZE_EXPAND_FILL
	left.add_theme_constant_override("separation",18)
	middle.add_child(left)
	var seats := HBoxContainer.new()
	seats.add_theme_constant_override("separation",18)
	left.add_child(seats)
	for i in range(2):
		var card := box(seats,Color("091727"))
		card.get_parent().size_flags_horizontal=Control.SIZE_EXPAND_FILL
		var top := HBoxContainer.new()
		card.add_child(top)
		var name_ := text("",22,WHITE,top)
		name_.size_flags_horizontal=Control.SIZE_EXPAND_FILL
		var you := text("YOU",13,MINT,top)
		var status := text("",14,WARN,card)
		status.autowrap_mode=TextServer.AUTOWRAP_WORD_SMART
		var deposits: Array = []
		for j in range(2):
			var row := HBoxContainer.new()
			card.add_child(row)
			var dname := text("",12,MUTED,row)
			dname.size_flags_horizontal=Control.SIZE_EXPAND_FILL
			var short := text("",12,MUTED,row)
			var progress := ProgressBar.new()
			progress.custom_minimum_size=Vector2(0,6)
			progress.show_percentage=false
			var track := StyleBoxFlat.new()
			track.bg_color=Color("172e46")
			track.set_corner_radius_all(3)
			progress.add_theme_stylebox_override("background",track)
			var fill := StyleBoxFlat.new()
			fill.bg_color=BLUE
			fill.set_corner_radius_all(3)
			progress.add_theme_stylebox_override("fill",fill)
			card.add_child(progress)
			deposits.append({"name":dname,"short":short,"bar":progress,"fill":fill})
		seat_cards.append({"card":card.get_parent(),"name":name_,"you":you,"status":status,"deposits":deposits})
	var next := box(left,Color("0b1d30"))
	text("NEXT STEP",12,MUTED,next)
	next_step=text("",24,WHITE,next)
	next_detail=text("",14,MUTED,next)
	next_detail.autowrap_mode=TextServer.AUTOWRAP_WORD_SMART
	error_label=text("",14,WARN,next)
	error_label.autowrap_mode=TextServer.AUTOWRAP_WORD_SMART

	var right := box(middle,Color("091727"))
	right.get_parent().custom_minimum_size=Vector2(400,0)
	text("TABLE TERMS",12,MUTED,right)
	for key in ["buyIn","bond","perSeat","pot","refund","until"]:
		var row := HBoxContainer.new()
		right.add_child(row)
		var caption := text({"buyIn":"Stake per seat","bond":"Admission bond","perSeat":"You pay in total","pot":"Winner's pot","refund":"Refund / bond lock","until":"Admission closes"}[key],14,MUTED,row)
		caption.size_flags_horizontal=Control.SIZE_EXPAND_FILL
		terms[key]=text("",14,WHITE,row)
	var note := text("The bond comes back after its lock. Payout fees apply to the pot, and the payout needs both players to sign; if it fails you reclaim your own stake after its lock. Every payment is approved in dcrpulse.",12,MUTED,right)
	note.autowrap_mode=TextServer.AUTOWRAP_WORD_SMART
	var spacer := Control.new()
	spacer.size_flags_vertical=Control.SIZE_EXPAND_FILL
	right.add_child(spacer)
	fund_button=action("Request stake",right,true)
	fund_button.pressed.connect(func(): fund_requested.emit())
	hide()

# update_lobby fills the screen from the backend's lobby description.
func update_lobby(lobby: Dictionary, fund_pending: bool) -> void:
	var match_id: String=lobby.get("match","")
	table_label.text="TABLE "+(match_id.substr(0,12)+"…" if match_id.length()>12 else match_id)+"   ·   "+String(lobby.get("stage",""))
	var rail: Array=lobby.get("steps",[])
	for i in range(steps.size()):
		var step: Dictionary=rail[i] if i<rail.size() else {"label":"","state":"todo"}
		var state: String=step.get("state","todo")
		step_states[i]=state
		var color := MINT if state=="done" else (WARN if state=="current" else MUTED)
		steps[i].bar.color=MINT if state=="done" else (WARN if state=="current" else TRACK)
		steps[i].name.text="%d  %s" % [i+1,step.get("label","")]
		steps[i].name.add_theme_color_override("font_color",color)
	var seats: Array=lobby.get("seats",[])
	for i in range(seat_cards.size()):
		var card: Dictionary=seat_cards[i]
		card.card.visible=i<seats.size()
		if i>=seats.size(): continue
		var seat: Dictionary=seats[i]
		card.name.text=seat.get("name","")
		card.you.visible=seat.get("you",false)
		card.status.text=seat.get("status","")
		card.status.add_theme_color_override("font_color",MINT if seat.get("status","")=="ready" else WARN)
		var deposits: Array=seat.get("deposits",[])
		for j in range(card.deposits.size()):
			var row: Dictionary=card.deposits[j]
			var d: Dictionary=deposits[j] if j<deposits.size() else {}
			row.name.text=d.get("name","")
			row.short.text=d.get("short","")
			row.bar.value=float(d.get("progress",0.0))*100.0
			row.fill.bg_color=WRONG if d.get("wrong",false) else (MINT if d.get("complete",false) else BLUE)
			row.short.add_theme_color_override("font_color",WRONG if d.get("wrong",false) else (MINT if d.get("complete",false) else MUTED))
	var t: Dictionary=lobby.get("terms",{})
	terms.buyIn.text="%s DCR" % t.get("buyIn","")
	terms.bond.text="%s DCR" % t.get("bond","")
	terms.perSeat.text="%s DCR" % t.get("perSeat","")
	terms.pot.text="%s DCR" % t.get("pot","")
	terms.refund.text="%d / %d blocks" % [int(t.get("refundBlocks",0)),int(t.get("bondBlocks",0))]
	terms.until.text="block %d" % int(t.get("until",0))
	next_step.text=lobby.get("nextStep","")
	next_detail.text=lobby.get("nextDetail","")
	if lobby.get("closed",false):
		error_label.text="Table closed: "+String(lobby.get("closedReason",""))
	elif lobby.get("stale",false):
		error_label.text="Evidence out of date. Reconnect before any payment."
	else:
		error_label.text=lobby.get("error","")
	fund_button.disabled=not lobby.get("canFund",false) or fund_pending
	fund_button.text="Stake requested · approve in dcrpulse" if fund_pending else "Request stake"
