using DevSmtp.Core.Models;

namespace DevSmtp.Core.Queries
{
    public class FindMessagesByEmail : IQuery<FindMessagesByEmailResult>
    {
        public FindMessagesByEmail(Email email)
        {
            this.Email = email;
        }

        public Email Email { get; }
    }
}
