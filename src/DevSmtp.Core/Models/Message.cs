namespace DevSmtp.Core.Models
{
    public class Message
    {
        public MessageId? Id { get; set; }
        public Email? From { get; set; }
        public IEnumerable<Email>? To { get; set; }
        public string? Data { get; set; }
    }
}
